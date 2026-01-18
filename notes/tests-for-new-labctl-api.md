## tests for new labctl api

| Command                                             | Expected                                                                              |
|-----------------------------------------------------|---------------------------------------------------------------------------------------|
| labctl get config proxmox                           | get config for ansible                                                                |
| labctl get config proxmox --for ansible             | get config for ansible                                                                |
| labctl get config proxmox --for ansible,ssh         | get config for ansible and ssh                                                        |
| labctl get config proxmox --for ssh,ansible         | get config for ansible and ssh                                                        |
| labctl get config proxmox --for ansible --for ssh   | get config for ansible and ssh                                                        |
| labctl get config proxmox --for ssh --for ansible   | get config for ansible and ssh                                                        |
| labctl check config proxmox                         | check config for ansible                                                              |
| labctl check config proxmox --for ansible           | check config for ansible                                                              |
| labctl check config proxmox --for ansible,ssh       | check config for ansible and ssh                                                      |
| labctl check config proxmox --for ssh,ansible       | check config for ansible and ssh                                                      |
| labctl check config proxmox --for ansible --for ssh | check config for ansible and ssh                                                      |
| labctl check config proxmox --for ssh --for ansible | check config for ansible and ssh                                                      |
| labctl set file proxmox                             | update ~/.ssh/config file with proxmox info                                           |
| labctl set file proxmox --for ssh                   | update ~/.ssh/config file with proxmox info                                           |
| labctl set file proxmox --for ansible               | error: invalid args: the following targets do not support writing to a file [ansible] |
