resource "scaleway_volume" "volume" {
  count = "1"
  name       = "${var.name}"
  size_in_gb = "${var.size_gb}"
  type       = "l_ssd"
}

resource "scaleway_volume_attachment" "attachment" {
  server = "${var.server_id}" 
  count = "1"
  volume = "${scaleway_volume.volume.id}"
}

