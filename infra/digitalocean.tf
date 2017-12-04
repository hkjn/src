provider "digitalocean" {
  token = "${chomp(file(var.digitalocean_token_file))}"
}

resource "digitalocean_ssh_key" "digitalocean1_id_rsa" {
  name       = "digitalocean1_id_rsa"
  public_key = "${var.digitalocean1_pubkey}"
}

resource "digitalocean_volume" "dropcore_disk0" {
  count = "${var.dropcore_enabled ? 1 : 0}"
  region      = "fra1"
  name        = "dropcore-disk0"
  size        = 170
  description = "Disk for dropcore"
}

resource "digitalocean_droplet" "dropcore" {
  count = "${var.dropcore_enabled ? 1 : 0}"
  image = "coreos-stable"
  name = "dropcore"
  region = "fra1"
  size = "2gb"
  ssh_keys = [
    "${digitalocean_ssh_key.digitalocean1_id_rsa.id}",
  ]
  volume_ids = [
    "${digitalocean_volume.dropcore_disk0.id}",
  ]
}


