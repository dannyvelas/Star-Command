package app

type resource string

const (
	ansiblePlaybookResource  resource = "ansible playbook"
	ansibleInventoryResource resource = "ansible inventory"
	terraformResource        resource = "terraform"
	sshResource              resource = "ssh"
)
