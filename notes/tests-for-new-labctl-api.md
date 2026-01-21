## tests for new labctl api

| Command                                             | Expected                                                                              | Result |
|-----------------------------------------------------|---------------------------------------------------------------------------------------|--------|
| labctl get config proxmox                           | get config for ansible                                                                | Passed |
| labctl get config proxmox --for ansible             | get config for ansible                                                                | Passed |
| labctl get config proxmox --for ssh                 | get config for ssh                                                                    | Passed |
| labctl get config proxmox --for ansible,ssh         | get config for ansible and ssh                                                        | Passed |
| labctl get config proxmox --for ssh,ansible         | get config for ansible and ssh                                                        | Passed |
| labctl get config proxmox --for ansible --for ssh   | get config for ansible and ssh                                                        | Passed |
| labctl get config proxmox --for ssh --for ansible   | get config for ansible and ssh                                                        | Passed |
| labctl check config proxmox                         | check config for ansible                                                              | Passed |
| labctl check config proxmox --for ansible           | check config for ansible                                                              | Passed |
| labctl check config proxmox --for ssh               | check config for ssh                                                                  | Passed |
| labctl check config proxmox --for ansible,ssh       | check config for ansible and ssh                                                      | Passed |
| labctl check config proxmox --for ssh,ansible       | check config for ansible and ssh                                                      | Passed |
| labctl check config proxmox --for ansible --for ssh | check config for ansible and ssh                                                      | Passed |
| labctl check config proxmox --for ssh --for ansible | check config for ansible and ssh                                                      | Passed |
| labctl set file proxmox                             | update ~/.ssh/config file with proxmox info                                           | Passed |
| labctl set file proxmox --for ssh                   | update ~/.ssh/config file with proxmox info                                           | Passed |
| labctl set file proxmox --for ansible               | error: invalid args: the following targets do not support writing to a file [ansible] | Passed |


| Command                                           | Expected                                                                              | Result |
|---------------------------------------------------|---------------------------------------------------------------------------------------|--------|
| GetConfig("proxmox", []string{"ansible"})         | json object which is config for ansible                                               | Passed |
| GetConfig("proxmox", []string{"ssh"})             | json object which is config for ssh                                                   | Passed |
| GetConfig("proxmox", []string{"ansible","ssh"})   | json object which has configs for both ssh and ansible                                | Passed |
| GetConfig("proxmox", []string{"ssh","ansible"})   | json object which has configs for both ssh and ansible                                | Passed |
| CheckConfig("proxmox", []string{"ansible"})       | table with diagnostic data for ansible                                                | Passed |
| CheckConfig("proxmox", []string{"ssh"})           | table with diagnostic data for ssh                                                    | Passed |
| CheckConfig("proxmox", []string{"ansible","ssh"}) | table with diagnostic data for ansible and ssh                                        | Passed |
| CheckConfig("proxmox", []string{"ssh","ansible"}) | table with diagnostic data for ansible and ssh                                        | Passed |
| SetFile("proxmox", []string{"ssh"})               | update ~/.ssh/config file with proxmox info                                           | Passed |
| SetFile("proxmox", []string{"ansible"})           | error: invalid args: the following targets do not support writing to a file [ansible] | Passed |
