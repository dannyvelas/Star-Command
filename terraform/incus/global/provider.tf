terraform {
  required_providers {
    incus = {
      source  = "lxc/incus"
      version = ">= 1.0.0"
    }
  }
  required_version = "~> 1.13.3"
}

provider "incus" {
  remote {
    name = var.incus_remote
  }
}
