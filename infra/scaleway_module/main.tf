resource "scaleway_ip" "hkjnprod" {
  count = "1"
}

resource "scaleway_server" "hkjnprod" {
  name           = "${var.machine_name}"
  image          = "${var.image}"
  type           = "C1"
  public_ip      = "${element(scaleway_ip.hkjnprod.*.ip, count.index)}"
}

module "volume1" {
  source = "./volume"
  name = "prodvolume1"
  server_id = "${scaleway_server.hkjnprod.id}"
}

module "volume2" {
  source = "./volume"
  name = "prodvolume2"
  server_id = "${scaleway_server.hkjnprod.id}"
}
