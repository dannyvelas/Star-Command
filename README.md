# Homelab setup

## Prerequisites
* A server
* A computer you can use to ssh into the server
* Terraform installed on that computer
* An ethernet connection between your router and the server

## Instructions
- Flash [Proxmox](https://www.proxmox.com/en/downloads) ISO onto a USB or SSD or disk and then connect that to your server so that you can boot your server with the Proxmox VE OS.
- In the setup, you'll be given an option to name your default proxmox node, pick a name but remember it because you will need it later.
- In setup, on the "Management Network Configuration" page:
  - For management interface, pick the network card that is being used for ethernet
  - For Hostname (FQDN), put proxmox.lan
  - For IP address (CIDR), pick an IP address that is not assigned by your router and that you can reserve for your server, let's suppose it's `1.2.3.4`
  - For gateway, put the IP address of your router
  - For DNS server, put either your router's IP address or use a public DNS server like `1.1.1.1` (Cloudflare) `8.8.8.8` (Google)
- Verify that from another computer you can `ping 1.2.3.4` and access `https://1.2.3.4:8006`.
- Add an `ssh` key so you have remote access to your server from another computer.
- `ssh` into your server and install sudo: `apt install sudo`. This is important because we will be creating a non-root user for yourself and a terraform user and both need `sudo` access.
- Create a non-root user for yourself:
  - Still in `ssh`, create a new user for yourself with full `sudo` permissions: `useradd -G sudo -m dannyvelasquez -s /bin/bash`
  - Create a password for that user: `passwd dannyvelasquez`. Save the password on Bitwarden.
  - Add `ssh` capability to `dannyvelasquez`:
    - Leave `ssh` and copy your SSH public key `cat /path/to/your/public/.ssh/key | pbcopy`
    - `ssh` into the server again and:
      - `cd /home/dannyvelasquez`
      - `mkdir .ssh`
      - `chmod 700 .ssh`
      - `chown dannyvelasquez:dannyvelasquez .ssh`
      - `touch authorized_keys`
      - Paste your key into `authorized_keys`
      - `chmod 600 authorized_keys`
      - `chown terraform:terraform authorized_keys`
    - Verify that it worked by running this from your local computer: `ssh -i /path/to/your/private/.ssh/key dannyvelasquez@1.2.3.4`
- Create an API token for Terraform. `ssh` into the VM and:
  - Create new user for terraform: `sudo pveum user add terraform@pve`
  - Create new role with terraform permissions: `sudo pveum role add Terraform -privs "Realm.AllocateUser, VM.PowerMgmt, VM.GuestAgent.Unrestricted, Sys.Console, Sys.Audit, Sys.AccessNetwork, VM.Config.Cloudinit, VM.Replicate, Pool.Allocate, SDN.Audit, Realm.Allocate, SDN.Use, Mapping.Modify, VM.Config.Memory, VM.GuestAgent.FileSystemMgmt, VM.Allocate, SDN.Allocate, VM.Console, VM.Clone, VM.Backup, Datastore.AllocateTemplate, VM.Snapshot, VM.Config.Network, Sys.Incoming, Sys.Modify, VM.Snapshot.Rollback, VM.Config.Disk, Datastore.Allocate, VM.Config.CPU, VM.Config.CDROM, Group.Allocate, Datastore.Audit, VM.Migrate, VM.GuestAgent.FileWrite, Mapping.Use, Datastore.AllocateSpace, Sys.Syslog, VM.Config.Options, Pool.Audit, User.Modify, VM.Config.HWType, VM.Audit, Sys.PowerMgmt, VM.GuestAgent.Audit, Mapping.Audit, VM.GuestAgent.FileRead, Permissions.Modify"`
  - Add role to previously created user: `sudo pveum aclmod / -user terraform@pve -role Terraform`
  - Create an API token for the user: `sudo pveum user token add terraform@pve provider --privsep=0`
  - Take note of the API token and save it into Bitwarden
- Create a user for Terraform with `sudo` privileges:
  - create terraform user: `useradd -m terraform`
  - give terraform user partial `sudo` permissions. Run: `visudo -f /etc/sudoers.d/terraform`. This will open a file. Put the following in that file:
    ```
    terraform ALL=(root) NOPASSWD: /sbin/pvesm
    terraform ALL=(root) NOPASSWD: /sbin/qm
    terraform ALL=(root) NOPASSWD: /usr/bin/tee /var/lib/vz/*
    ```
  - Add `ssh` capability to `terraform`:
    - Leave `ssh` and copy your SSH public key `cat /path/to/your/public/.ssh/key | pbcopy`
    - `ssh` into the server again and:
      - `cd /home/terraform`
      - `mkdir .ssh`
      - `chmod 700 .ssh`
      - `chown terraform:terraform .ssh`
      - `touch authorized_keys`
      - Paste your key into `authorized_keys`
      - `chmod 600 authorized_keys`
      - `chown terraform:terraform authorized_keys`
    - Verify that it worked by running this from your local computer: `ssh -i /path/to/your/private/.ssh/key terraform@1.2.3.4 sudo pvesm apiinfo`. You should see something like this:
      ```
      APIVER 12
      APIAGE 3
      ```
- Create a file in this directory called `terraform.tfvars` it should look like this:
```
endpoint        = "https://1.2.3.4:8006/"
api_token       = "terraform@pve!provider=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
ssh_public_key  = "/path/to/your/public/.ssh/key"
ssh_private_key = "/path/to/your/private/.ssh/key"
node            = "whatever-node-name-you-chose-earlier"
```
- `ssh` into your server and create the `/mnt/media` directory
- Run `terraform apply`. This should create an Ubuntu VM that has a shared mount to the `/mnt/media` directory of its host.
