# Internals

## How `iac setup` works

`iac setup` provisions one or more hosts end-to-end. It runs the same sequence of steps for each host, detecting existing cluster state along the way so that subsequent hosts join rather than reinitialize.

### 1. Generate Ansible inventory

An Ansible inventory file is generated from `iac.yml`. It includes all configured hosts and their VMs. VMs are configured with `ProxyJump` pointing to their parent host, so that later Ansible steps can reach them through the host without requiring the VMs to be directly accessible on the network.

### 2. Bootstrap hosts

```
ansible-playbook bootstrap-server.yml -u root
```

Runs against all hosts. Hardens the OS: configures UFW, enforces SSH key-only authentication, and enables unattended security updates.

### 3. Set up hosts

```
ansible-playbook setup-host.yml
```

Runs against all hosts. Installs and configures host-level services: Incus (hypervisor), WireGuard (on the designated VPN host), and Traefik (reverse proxy).

---

The following steps run once per host, in order:

### 4. Register host in `~/.ssh/config`

The host is added to `~/.ssh/config` on your workstation if not already present, so that subsequent SSH and Ansible commands can refer to it by name.

### 5. Join OVN overlay network

OVN provides encrypted east-west networking between hosts.

- **First host:** initializes the OVN central database
- **Subsequent hosts:** join as a chassis node

### 6. Join k3s cluster (server)

k3s runs on each host as a server node.

- **First host:** initializes the cluster
- **Subsequent hosts:** join the existing cluster as additional server nodes

### 7. Create VMs with Terraform

```
terraform init && terraform apply
```

Terraform provisions the VMs for this host via Incus. Each VM is on a private NAT subnet and is not directly reachable from the physical network.

### 8. Bootstrap VMs

```
ansible-playbook bootstrap-server.yml --limit vms -u root
```

Same OS hardening as step 2, now applied to the newly created VMs.

### 9. Set up VMs

```
ansible-playbook setup-vm.yml
```

Runs against all VMs for this host. Installs VM-level services (Docker, storage mount points).

### 10. Register VMs in `~/.ssh/config`

Each VM is added to `~/.ssh/config` with a `ProxyJump` directive pointing to its parent host, making it reachable by name from your workstation.

### 11. Join k3s cluster (agent)

Each VM joins the k3s cluster as a worker node, making it available for workload scheduling.

### 12. Join Incus cluster

- **First host:** initializes the Incus cluster and sets it as the active remote
- **Subsequent hosts:** join the existing Incus cluster

---

Once all hosts are processed, `iac setup` prints a summary of the Incus cluster with `incus list`.
