# Homelab setup

## Prerequisites
* A server
* A computer you can use to ssh into the server
* Terraform installed on that computer
* An ethernet connection between your router and the server

## Instructions
1. Flash [Proxmox](https://www.proxmox.com/en/downloads) ISO onto a USB or SSD or disk and then connect that to your server so that you can boot your server with the Proxmox VE OS.
2. In setup, on the "Management Network Configuration" page:
  - For management interface, pick the network card that is being used for ethernet
  - For Hostname (FQDN), put proxmox.lan
  - For IP address (CIDR), pick an IP address that is not assigned by your router and that you can reserve for your server
  - For gateway, put the IP address of your router, let's suppose it's `1.2.3.4`
  - For DNS server, put either your router's IP address or use a public DNS server like `1.1.1.1` (Cloudflare) `8.8.8.8` (Google)
3. Verify that from another computer you can `ping 1.2.3.4` and access `https://1.2.3.4:8006`.
4. Add an `ssh` key so you can access a pseudo-terminal of your server from another computer.
5. Create an API token for Terraform. `ssh` into the VM and:
  - Create new user for terraform: `sudo pveum user add terraform@pve`
  - Create new role with terraform permissions: `sudo pveum role add Terraform -privs "Realm.AllocateUser, VM.PowerMgmt, VM.GuestAgent.Unrestricted, Sys.Console, Sys.Audit, Sys.AccessNetwork, VM.Config.Cloudinit, VM.Replicate, Pool.Allocate, SDN.Audit, Realm.Allocate, SDN.Use, Mapping.Modify, VM.Config.Memory, VM.GuestAgent.FileSystemMgmt, VM.Allocate, SDN.Allocate, VM.Console, VM.Clone, VM.Backup, Datastore.AllocateTemplate, VM.Snapshot, VM.Config.Network, Sys.Incoming, Sys.Modify, VM.Snapshot.Rollback, VM.Config.Disk, Datastore.Allocate, VM.Config.CPU, VM.Config.CDROM, Group.Allocate, Datastore.Audit, VM.Migrate, VM.GuestAgent.FileWrite, Mapping.Use, Datastore.AllocateSpace, Sys.Syslog, VM.Config.Options, Pool.Audit, User.Modify, VM.Config.HWType, VM.Audit, Sys.PowerMgmt, VM.GuestAgent.Audit, Mapping.Audit, VM.GuestAgent.FileRead, Permissions.Modify"`
  - Add role to previously created user: `sudo pveum aclmod / -user terraform@pve -role Terraform`
  - Create an API token for the user: `sudo pveum user token add terraform@pve provider --privsep=0`
  - Take note of the API token and save it into bitwarden
6. `ssh` into your server and create the `/mnt/media` directory
7. Run `terraform apply`. This should create an Ubuntu VM that has a shared mount to the `/mnt/media` directory of its host.
