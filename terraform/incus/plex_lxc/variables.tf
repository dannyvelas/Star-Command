variable "incus_remote" {
  type    = string
  default = "local"
}

variable "incus_project" {
  description = "The Incus project"
  type        = string
  default     = "default"
}

variable "ip" {
  description = "Static IP for the Plex container"
  type        = string
}

variable "ssh_public_key_path" {
  type    = string
  default = "~/.ssh/id_ed25519.pub"
}
