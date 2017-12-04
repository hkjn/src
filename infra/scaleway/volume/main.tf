resource "scaleway_volume" "volume" {
  name       = "${var.name}"
  count = "${var.enabled ? 1 : 0}"
  size_in_gb = "${var.size_gb}"
  type       = "l_ssd"
}

resource "scaleway_volume_attachment" "attachment" {
  server = "${var.server_id}" 
  count = "${var.enabled ? 1 : 0}"
  volume = "${scaleway_volume.volume.id}"
}

