
| Function       | arg1     | arg2      |
|----------------|----------|-----------|
| AnsibleRun     | playbook | preflight |
| SSHAdd         | <host>   | preflight |
| TerraformApply |          | preflight |
| InventoryAdd   | <host>[] | preflight |
| Setup          | <host>[] | preflight |


| Resource  | verb             | arg2     | arg3      |
|-----------|------------------|----------|-----------|
| Ansible   | bootstrap-server |          | preflight |
| Ansible   | setup-host       |          | preflight |
| Ansible   | setup-remote     |          | preflight |
| Ansible   | setup-vm         |          | preflight |
| SSH       | add              | <host>   | preflight |
| Terraform | apply            |          | preflight |
| Inventory | add              | <host>[] | preflight |
| Setup     |                  | <host>[] | preflight |
