terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.83.2"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5.3"
    }
  }
  required_version = "~> 1.13.3"
}

provider "proxmox" {
  endpoint = var.endpoint
  username = var.username
  password = var.password
  insecure = true
}

data "local_file" "ssh_public_key" {
  filename = var.ssh_public_key
}

resource "proxmox_virtual_environment_download_file" "ubuntu_lxc_template" {
  content_type = "vztmpl"
  datastore_id = "local"
  node_name    = var.node
  url          = "http://download.proxmox.com/images/system/ubuntu-24.04-standard_24.04-2_amd64.tar.zst"
  file_name    = "ubuntu-24.04-standard_24.04-2_amd64.tar.zst"
}

resource "proxmox_virtual_environment_container" "plex_lxc" {
  description  = "Managed by Terraform"
  tags         = ["terraform", "ubuntu"]
  node_name    = var.node
  vm_id        = 100
  unprivileged = true

  cpu {
    cores = 2
  }

  memory {
    dedicated = 1024
  }

  disk {
    datastore_id = "local-lvm"
    size         = 20
  }

  network_interface {
    name     = "eth0"
    bridge   = "vmbr0"
    firewall = true
  }

  # this initialization block works because:
  # the `proxmox_virtual_environment_download_file` resource created a VM template and
  # stored it in the proxmox "local" storage. this template is configured so that when
  # a VM using this template boots for the first time, there will be a special
  # cloud-init drive in it. this allows us to pass data into the VM like SSH keys, hostname,
  # network config, etc
  initialization {
    hostname = "terraform-provider-proxmox-ubuntu-container"

    ip_config {
      ipv4 {
        address = "${var.ip}/24"
        gateway = var.router_ip
      }
    }

    user_account {
      keys = [
        trimspace(data.local_file.ssh_public_key.content)
      ]
    }

    dns {
      servers = ["1.1.1.1"]
    }
  }

  operating_system {
    template_file_id = proxmox_virtual_environment_download_file.ubuntu_lxc_template.id
    type             = "ubuntu"
  }

  mount_point {
    # bind mount for media
    volume = "/mnt/media"
    path   = "/mnt/media"
  }

  mount_point {
    # bind mount for plex metadata
    volume = "/mnt/media/plex-config"
    path   = "/var/lib/plexmediaserver/Library/Application Support/Plex Media Server"
  }
}

# This resource waits for the LXC to be ready, then reaches through the host to flip the SSH port
# from 22 to 17031
resource "terraform_data" "bootstrap_ssh" {
  # This replaces "depends_on". If the VM ID changes, this re-runs.
  triggers_replace = [
    proxmox_virtual_environment_container.plex_lxc.id
  ]

  connection {
    type        = "ssh"
    user        = "admin"
    port        = 17031
    host        = var.host_ip
    private_key = file(var.ssh_private_key)
  }

  provisioner "remote-exec" {
    inline = [
      "timeout 60s bash -c 'until sudo /usr/sbin/pct exec 100 -- ls /etc/ssh/sshd_config; do sleep 2; done'",
      "sudo /usr/sbin/pct exec 100 -- mkdir -p /etc/systemd/system/ssh.socket.d",
      "sudo /usr/sbin/pct exec 100 -- bash -c \"echo -e '[Socket]\nListenStream=\nListenStream=17031' > /etc/systemd/system/ssh.socket.d/listen.conf\"",
      "sudo /usr/sbin/pct exec 100 -- systemctl daemon-reload",
      "sudo /usr/sbin/pct exec 100 -- systemctl stop ssh.socket",
      "sudo /usr/sbin/pct exec 100 -- systemctl start ssh.socket",
      "sudo /usr/sbin/pct exec 100 -- sed -i 's/^#?Port 22/Port 17031/' /etc/ssh/sshd_config"
    ]
  }
}

# enable firewall options at container level which is default deny
resource "proxmox_virtual_environment_firewall_options" "plex_fw_options" {
  node_name    = var.node
  container_id = proxmox_virtual_environment_container.plex_lxc.vm_id

  enabled       = true
  input_policy  = "DROP"   # Everything not allowed is blocked
  output_policy = "ACCEPT" # Allow the container to reach out for updates
}

# enable firewall rule which associates the "plex" security group with this plex_lxc
# the "plex" security group should have been created prior (it is defined in terraform/global)
resource "proxmox_virtual_environment_firewall_rules" "plex_rules" {
  node_name    = var.node
  container_id = proxmox_virtual_environment_container.plex_lxc.vm_id

  rule {
    security_group = "plex"
    iface          = "net0"
  }

  rule {
    security_group = "guest_mgmt"
    iface          = "net0"
  }
}
