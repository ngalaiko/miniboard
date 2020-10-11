provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_region" "ams3" {
  slug = "ams3"
}

data "digitalocean_sizes" "main" {
  filter {
    key    = "slug"
    values = ["s-1vcpu-1gb"]
  }
}

module "backend" {
  source = "./backend"

  binary_path = var.backend_binary_path
  region      = data.digitalocean_region.ams3.slug
  size        = element(data.digitalocean_sizes.main.sizes, 0).slug
  ip_whitelist = [
    "155.4.221.5/32",
  ]
}
