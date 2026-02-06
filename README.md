# Homelab infra and playbooks

## Prerequisites
* A server (Ubuntu/Debian) connected to Ethernet.
* A computer you can use to ssh into the server (Controller).
* [Terraform](https://developer.hashicorp.com/terraform/install) installed.
* [Ansible](https://formulae.brew.sh/formula/ansible) installed.
* A [Tailscale](https://login.tailscale.com/start) account.
* SMTP credentials (for update notifications).

<details>

<summary><h2>1. Set Ansible variables</h2></summary>

- Pick an admin password for your home server. Save it into Bitwarden as well.
- Run: `ansible-vault create ./ansible/host_vars/proxmox_server/vault.yml`, using a vault password. Save this vault password in Bitwarden.
- In the content of that file put:
  ```
  vault_admin_password: "<admin password for your home server>"
  ```
- Update `./ansible/inventory.ini` so that the `proxmox` host has IP address `1.2.3.4`.
- Run `ansible-vault create ./ansible/group_vars/all/vault.yml`, and add the following:
  ```
  smtp_user: "your-email@example.com"
  smtp_pass: "your 16-character code if gmail, otherwise regular password"
  ```
- Make sure `./ansible/group_vars/all/all.yml` looks something like this:
  ```
  ssh_port: 17031
  ansible_user: admin
  ansible_port: "{{ ssh_port }}"
  ```

</details>

<details>

<summary><h2>2. Bootstrap Server (Metal)</h2></summary>

- Flash Ubuntu Server or Debian onto your server.
- Configure SSH access.
- **Install Ansible Requirements**:
  ```bash
  ansible-galaxy install -r ansible/requirements.yml
  ```
- **Bootstrap Incus & Harden**:
  This sets up Incus, initializes storage/network, hardens SSH (Port 17031), and sets up auto-updates.
  ```bash
  ansible-playbook -i ansible/inventory.ini ansible/setup-incus.yml --ask-become-pass --ask-vault-pass
  ```
  *(Note: Ensure your `inventory.ini` has the correct IP for the `incus` host entry)*

</details>

<details>

<summary><h2>3. Configure Incus Remote</h2></summary>

- On your laptop, authorize the Incus remote so Terraform can talk to it:
  ```bash
  incus remote add my-homelab <Server-IP>
  incus remote switch my-homelab
  ```
- Verify connection:
  ```bash
  incus list
  ```

</details>


<details>

<summary><h2>4. Global Infrastructure (Terraform)</h2></summary>

- `cd terraform/incus/global`
- Create `terraform.tfvars` if needed (defaults are usually fine):
  ```hcl
  incus_remote = "my-homelab"
  ```
- Run `terraform init`.
- Run `terraform apply`.
- This creates the `basic` (network/disk) and `management` (ssh keys) profiles used by all instances.

</details>

<details>

<summary><h2>5. Deploy Plex (LXC)</h2></summary>

- `cd terraform/incus/plex_lxc`
- Create `terraform.tfvars`:
  ```hcl
  incus_remote = "my-homelab"
  ip           = "10.0.100.85"
  # Optional: ssh_public_key_path = "~/.ssh/id_ed25519.pub"
  ```
- Run `terraform apply`.
- This creates the Plex container, binds `/mnt/media`, and exposes port `32400`.
- **Post-Provisioning**:
  Run the Ansible setup to install the Plex binary and configure ownership:
  ```bash
  ansible-playbook -i ansible/inventory.ini ansible/setup-plex-lxc.yml --ask-vault-pass
  ```

</details>

<details>

<summary><h2>6. Deploy WireGuard (VM)</h2></summary>

- `cd terraform/incus/wireguard_vm`
- Create `terraform.tfvars`:
  ```hcl
  incus_remote = "my-homelab"
  ip           = "10.0.100.84"
  # ssh_public_key_path = "..."
  ```
- Run `terraform apply`.
- This creates the VM with UDP port forwarding (`51820`) and SSH mapping (`17031` -> `22`).

</details>
