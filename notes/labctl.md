# labctl spec

An internal CLI to configure a homelab

## labctl grammar v1

| Action    | Host-Alias | Flags     |
|-----------|------------|-----------|
| resolve   | proxmox    |           |
| resolve   | proxmox    | --dry-run |
| setup-ssh | proxmox    |           |

## labctl grammar v2

| Action | Resource | Host-Alias |
|--------|----------|------------|
| get    | config   | proxmox    |
| set    | ssh      | proxmox    |
| check  | reqs     | proxmox    |

## labctl grammar v3

| Action | Resource | Host-Alias | Flags             |
|--------|----------|------------|-------------------|
| get    | config   | proxmox    | --for ansible     |
| set    | file     | proxmox    | --for ssh         |
| check  | config   | proxmox    | --for ansible,ssh |
