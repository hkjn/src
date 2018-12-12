provider "scaleway" {
	organization = "${var.scaleway_organization}"
	token = "${var.scaleway_token}"
	region = "par1"
}

resource "scaleway_ip" "lab" {
	count = "1"
}


resource "scaleway_server" "lab" {
	name           = "lab"
	count          = 0
	image          = "${var.image}"
	type           = "C1"
	public_ip      = "${element(scaleway_ip.lab.*.ip, count.index)}"
}
