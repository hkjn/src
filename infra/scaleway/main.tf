resource "scaleway_ip" "hkjnprod" {
  count = "${var.enabled ? 1 : 0}"
}

resource "scaleway_server" "hkjnprod" {
  count = "${var.enabled ? 1 : 0}"
  name           = "prod.hkjn.me"
  image          = "${var.image}"
  type           = "C1"
  public_ip      = "${element(scaleway_ip.hkjnprod.*.ip, count.index)}"
}

module "volume1" {
  source = "./volume"
  name = "prodvolume1"
  server_id = "${scaleway_server.hkjnprod.id}"
  enabled = "${var.enabled}"
}

module "volume2" {
  source = "./volume"
  name = "prodvolume2"
  server_id = "${scaleway_server.hkjnprod.id}"
  enabled = "${var.enabled}"
}

module "volume3" {
  source = "./volume"
  name = "prodvolume3"
  server_id = "${scaleway_server.hkjnprod.id}"
  enabled = "${var.enabled}"
}

module "volume4" {
  source = "./volume"
  name = "prodvolume4"
  server_id = "${scaleway_server.hkjnprod.id}"
  enabled = "${var.enabled}"
}
