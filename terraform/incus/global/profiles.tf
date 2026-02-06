locals {
  ssh_key = trimspace(file(var.ssh_public_key_path))
}

resource "incus_profile" "basic" {
  name = "basic"
  description = "Basic networking and root disk configuration"

  device {
    name = "eth0"
    type = "nic"
    properties = {
      network = var.network_bridge
      name    = "eth0"
    }
  }

  device {
    name = "root"
    type = "disk"
    properties = {
      path = "/"
      pool = "default"
    }
  }
}

resource "incus_profile" "management" {
  name = "management"
  description = "Management access (SSH keys)"

  config = {
    "user.user-data" = <<-EOT
      #cloud-config
      ssh_authorized_keys:
        - ${local.ssh_key}
      packages:
        - curl
        - git
        - vim
        - htop
    EOT
  }
}
