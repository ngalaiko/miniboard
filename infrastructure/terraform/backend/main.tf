resource "tls_private_key" "backend" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "digitalocean_ssh_key" "backend" {
  name       = "backend"
  public_key = tls_private_key.backend.public_key_openssh
}


data "digitalocean_images" "ubuntu" {
  filter {
    key    = "slug"
    values = ["ubuntu-20-04-x64"]
  }
  filter {
    key    = "regions"
    values = [var.region]
  }
  sort {
    key       = "created"
    direction = "desc"
  }
}

resource "digitalocean_droplet" "backend" {
  count = var.replicas

  image              = element(data.digitalocean_images.ubuntu.images, 0).slug
  name               = "backend-${count.index}"
  region             = var.region
  size               = var.size
  monitoring         = true
  private_networking = true
  tags               = ["backend"]

  ssh_keys = [
    digitalocean_ssh_key.backend.fingerprint,
  ]

  connection {
    host        = digitalocean_droplet.backend[count.index].ipv4_address
    user        = "root"
    type        = "ssh"
    private_key = tls_private_key.backend.private_key_pem
    timeout     = "2m"
  }

  provisioner "file" {
    source      = var.binary_path
    destination = "/root/miniboard"
  }

  user_data = <<EOT
  #!/bin/bash
  
  echo Binary SHA256:${filesha256(var.binary_path)}
  EOT

  lifecycle {
    create_before_destroy = true
  }
}

resource "digitalocean_firewall" "backend" {
  name = "backend"

  tags = ["backend"]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = var.ip_whitelist
  }
}
