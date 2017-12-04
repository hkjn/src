provider "scaleway" {
  organization = "${chomp(file(var.scaleway_organization_file))}"
  token        = "${chomp(file(var.scaleway_token_file))}"
  region = "${var.scaleway_region}"
}

module "scaleway" {
  source = "./scaleway"
  enabled = "${var.hkjnprod_enabled}"
}
