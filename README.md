# Homelab infra and playbooks

Infrastructure as Code for a homelab environment. A single Go CLI (`labctl`) reads your configurations and orchestrates Terraform and Ansible to provision one or more servers with a hypervisor, a WireGuard VPN, a Traefik reverse proxy, OVN networking, and a k3s cluster. Engineers then deploy any dockerized service with `kubectl` — no application-specific code in the platform itself.

## Architecture

```
Home LAN (192.168.1.0/24)
  |
  +-- Home Gateway/Router
  |     \-- UDP 51820 forwarded -> Host 01
  |         TCP 443 forwarded -> Host 01
  |
  +-- Host 01 (192.168.1.10)
  |     +-- WireGuard (wg0, VPN endpoint)
  |     +-- UFW (NAT, port forwarding, egress blocking)
  |     +-- OVN (overlay networking between hosts)
  |     +-- k3s server (workload scheduler)
  |     \-- VM (192.168.122.50) <- private subnet, not on LAN
  |           +-- k3s agent
  |           +-- Traefik (reverse proxy, subdomain routing)
  |           \-- workload containers (read-only root)
  |
  +-- Host 02 (192.168.1.11)
  |     +-- UFW (NAT, port forwarding, egress blocking)
  |     +-- OVN (overlay networking between hosts)
  |     +-- k3s server
  |     \-- VM (192.168.123.50) <- private subnet, not on LAN
  |           +-- k3s agent
  |           +-- Traefik (reverse proxy, subdomain routing)
  |           \-- workload containers (read-only root)
  |
  +-- ...more hosts as needed...
```

**Two-layer isolation**: Trusted infrastructure (WireGuard, KVM, OVN) runs on the host. Application workloads run inside VMs behind NAT. A container escape lands in the VM kernel, not the host.

### Security layers

1. **OS hardening** — UFW firewall, SSH key-only auth, unattended security updates
2. **NAT networking** — VMs on private subnets, invisible to the home LAN
3. **Egress blocking** — VMs cannot initiate connections to other LAN devices
4. **VM kernel isolation** — Container breakouts are contained by the VM boundary
5. **OVN overlay** — Encrypted east-west traffic between hosts, isolated from the physical network
6. **Reverse proxy** — Single ingress point with TLS, security headers, rate limiting
7. **Read-only containers** — Immutable root filesystem, writable only for data volumes

### Monitoring and alerting

- **Grafana dashboard** — real-time metrics for all deployed services (resource usage, uptime, response times)
- **Alerting** — automatic notifications when any service goes down or degrades (email, Slack, PagerDuty)

### Go links

Internal short URLs for quick access to services and dashboards:

| Link | Destination |
|------|-------------|
| `go/grafana` | Monitoring dashboard |
| `go/alerts` | Alert configuration |
| `go/vpn` | VPN client setup guide |

## Prerequisites

### Hardware

- One or more Debian 12 servers, each with:
  - 8 GB+ RAM
  - Ethernet to home LAN

### Software (on your workstation)

- [Go](https://go.dev/) 1.21+ (to build the CLI)
- [Terraform](https://www.terraform.io/) (latest stable)
- [Ansible](https://docs.ansible.com/) (latest stable)
- [WireGuard client](https://www.wireguard.com/install/) (for VPN access)
- A C toolchain (e.g., `xcode-select --install` on macOS)
- SSH client with key-based auth configured to all servers

### Network preparation

Forward UDP port **51820** and TCP port **443** on your home gateway/router to the server that will act as the WireGuard endpoint and ingress point.

### Remote access (optional)

If you want to access services from outside your home LAN without VPN, set up a DDNS hostname before provisioning:

1. Create an account with a DDNS provider (e.g., DuckDNS, No-IP)
2. Register a hostname (e.g., `home.example.com`)
3. Note your login credentials — you'll add them to `config/`

During provisioning, the system automatically installs ddclient to keep the hostname pointed at your public IP. If you skip this step, everything still works on your home LAN and over VPN.

## Setup

### 1. Clone and build

```bash
git clone <repo-url>
cd homelab
make
```

### 2. Configure

Fill in your values:

```bash
bws secret create <KEY> <VALUE> <PROJECT_ID> # Bitwarden Secrets Manager : priority #1
vim ./config/<HOST_NAME>.yml                 # Host-specific config file : priority #2
vim ./config/all.yml                         # Global config file        : priority #3
cp .env.example .env                         # Environment variables     : priority #4
```

The `labctl` CLI reads these configuration values and generates all tool-specific configs (Ansible inventory, Terraform tfvars) into `.generated/`. You never edit those files directly.

## Usage

### Bootstrap infrastructure

Provision each server — this hardens the OS, installs WireGuard (on the designated VPN host), installs a hypervisor (Incus), creates a workload VM with Traefik, configures OVN overlay networking, and joins k3s:

```bash
labctl setup                     # setup all hosts specified in your config
labctl setup --host <your-host>  # setup only one host. If a cluster already exists, join this host to that cluster. Otherwise, initialize a new cluster
```

### Generate VPN client configs

```bash
labctl generate vpn-client --name "danny-laptop"
labctl generate vpn-client --name "danny-phone"
```

Client configs are saved to `.generated/vpn-clients/`. Import them into the WireGuard app on each device, then delete the `.conf` files from your workstation — they contain the client's private key and preshared key. The `.generated/` directory is gitignored and the files are created with `0600` permissions, but they should be treated as sensitive and not kept around longer than needed.

### Set up local DNS

Deploy CoreDNS on the cluster so that `*.home.example.com` resolves to the ingress host's LAN IP. This is what allows devices on your LAN to reach services like `plex.home.example.com`.

```bash
kubectl apply -f services/coredns.yml
```

Then update your router's DHCP settings to use the cluster as the DNS server. This is a one-time manual change — after this, every device on the LAN resolves service subdomains automatically.

### Deploy services

This platform is application-agnostic. After infrastructure is set up, deploy any service using standard Kubernetes manifests and kubectl. See `services/` for full working examples.

k3s handles scheduling and restarts. Traefik picks up new services automatically via Ingress resources and routes by subdomain.

```bash
kubectl apply -f services/plex.yml
kubectl apply -f services/sonarr.yml
kubectl apply -f services/radarr.yml
kubectl apply -f services/bazarr.yml
kubectl apply -f services/grafana.yml
kubectl apply -f services/golinks.yml
```

### Verify

```bash
labctl status # cluster overview: hosts, services, VPN, k3s
```

### Tear down

```bash
labctl teardown # destroy all VMs via Terraform
```

## CLI reference

```
labctl <command> [options]

Commands:
  setup                Setup a physical host (hardening, hypervisor, VPN, Reverse Proxy, VM, OVN, k3s)
  generate vpn-client  Generate a WireGuard client config
  status               Show cluster status (hosts, services, VPN, k3s)
  teardown             Tear down all VMs
  version              Print version
```

## Project structure

```
config/                      # config directory
  all.yml                    # global config
  host-01.yml                # host-specific config
services/                    # service manifests (one per app)
cmd/                         # Go CLI source
  labctl/                    # entrypoint and command routing
ansible/
  playbooks/                 # ansible playbooks
  roles/                     # ansible roles
terraform/                   # incus VM lifecycle
.generated/                  # auto-generated configs (gitignored)
```

## Scaling to new hardware

When you add or migrate servers:

1. Update `config/` with the new host IPs and storage paths
2. Run `labctl setup --host <new-host>` for each new host
3. k3s automatically joins the new node to the cluster
4. OVN extends the overlay network to the new host
5. Deploy services with `kubectl` — no changes to manifests needed, k3s schedules across the cluster

## Tech stack

| Component    | Tool                  | Why                                          |
|--------------|-----------------------|----------------------------------------------|
| CLI          | Go                    | Single binary, no runtime dependencies       |
| Provisioning | Ansible               | Agentless, SSH-based, idempotent             |
| VM lifecycle | Terraform + libvirt   | Declarative, reproducible                    |
| Scheduling   | k3s                   | Lightweight Kubernetes, rolling updates, auto-restart |
| Networking   | OVN                   | Overlay networking, encrypted east-west traffic |
| Firewall     | UFW                   | Simple, readable firewall rules              |
| VPN          | WireGuard             | Kernel-level, no third-party trust           |
| Proxy        | Traefik               | Auto-discovery, TLS, subdomain routing       |
| Monitoring   | Grafana + Prometheus  | Dashboards, alerting, service health         |
| Go links     | golinks               | Internal short URLs for quick service access |
| OS           | Debian 12             | Stable, security updates, KVM support        |
