terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }
    postgresql = {
      source = "terraform-providers/postgresql"
    }
  }
  required_version = ">= 0.13"
}
