# Configure the Google Cloud provider
provider "google" {
  credentials = "${file(var.gcloud_credentials)}"
  project     = "${var.gcloud_project}"
  region      = "${var.gcloud_region}"
}

#
# Network and firewalls.
#

resource "google_compute_network" "default" {
  name                    = "tf-net0"
  auto_create_subnetworks = "true"
}

resource "google_compute_firewall" "allow_ssh" {
  name    = "tf-allow-ssh-ping"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["22", "6200"]
  }
  allow {
    protocol = "udp"
    ports    = ["60000-60100"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["ssh"]
}

resource "google_compute_firewall" "http" {
  name = "tf-allow-http"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "tcp"
    ports = ["80"]
  }
  allow {
    protocol = "tcp"
    ports = ["443"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["http"]
}

resource "google_compute_firewall" "bitcoin" {
  name = "tf-allow-bitcoin"
  network = "${google_compute_network.default.name}"
  allow {
    protocol = "tcp"
    ports = ["8333"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags = ["bitcoin"]
}

#
# Storage
#

#
# decenter.world persistent disk
#
resource "google_compute_disk" "zdisk1" {
  name  = "zdisk1"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
  size  = "20"
}

#
# admin2.hkjn.me persistent disk
#
resource "google_compute_disk" "admin2-disk" {
  name  = "admin2-disk"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
  size  = "50"
}

#
# builder.hkjn.me persistent disk
#
resource "google_compute_disk" "builder_disk" {
  count        = "${var.builder_enabled ? 1 : 0}"
  name  = "builder-disk"
  type  = "pd-ssd"
  zone  = "europe-west3-b"
  size  = "200"
}

# GCP instances
#
# TODO: Set sshd port in bootstrap:
# cat /etc/systemd/system/sshd.socket.d/10-sshd-listen-ports.conf
# [Socket]
# ListenStream=
# ListenStream=6200

resource "google_compute_instance" "admin2" {
  count = 1
  name         = "admin2"
  description  = "admin2.hkjn.me server"
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["ssh", "http"]
  disk {
    image = "${var.ubuntu_image}"
  }
  attached_disk {
    source = "${google_compute_disk.admin2-disk.self_link}"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "zero:${var.admin2_pubkey}"
    version = "${var.infra_version}"
  }
}

resource "google_compute_instance" "decenter-world" {
  count = 1
  name         = "decenter-world"
  description  = "decenter.world server"
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["ssh", "http"]
  disk {
    image = "${var.coreos_alpha_image}"
  }
  attached_disk {
    source = "${google_compute_disk.zdisk1.self_link}"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.admin2_decenter_world_pubkey}"
    version = "${var.infra_version}"
    user-data = "${file("bootstrap/decenter_world.json")}"
  }
}

resource "google_compute_instance" "builder" {
  count        = "${var.builder_enabled ? 1 : 0}"
  name         = "builder"
  description  = "Temporary builder server."
  machine_type = "n1-standard-8"
  zone         = "europe-west3-b"
  tags         = ["ssh", "http"]
  disk {
    image = "${var.coreos_alpha_image}"
  }
  attached_disk {
    source = "${google_compute_disk.builder_disk.self_link}"
  }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.admin2_builder_pubkey}"
    version = "${var.infra_version}"
    user-data = "${file("bootstrap/builder.json")}"
  }
}

resource "google_compute_instance" "guac" {
  count = "${var.guac_enabled ? 1 : 0}"
  name         = "guac"
  description  = "The server for guac.hkjn.me."
  machine_type = "f1-micro"
  zone         = "europe-west3-a"
  tags = ["ssh", "http"]
  disk { image = "${var.coreos_alpha_image}" }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.admin2_guac_pubkey}"
  }
}

resource "google_compute_instance" "blockpress_me" {
  count = "${var.blockpress_me_enabled ? 1 : 0}"
  name         = "blockpress-me"
  description  = "The server for blockpress.me."
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["ssh", "http"]
  disk { image = "${var.coreos_alpha_image}" }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.admin1_blockpress_me_pubkey}"
  }
}

resource "google_compute_instance" "elentari-world" {
  count = "${var.elentari_world_enabled ? 1 : 0}"
  name         = "elentari-world"
  description  = "The server for elentari.world."
  machine_type = "f1-micro"
  zone         = "europe-west3-b"
  tags = ["ssh", "http"]
  disk { image = "${var.coreos_alpha_image}" }
  network_interface {
    network = "${google_compute_network.default.name}"
    access_config {} # Ephemeral IP
  }
  metadata {
    sshKeys = "core:${var.aruna_pubkey_arunallave}"
  }
}
