provider "scaleway" {
  organization = "${chomp(file(var.scaleway_organization_file))}"
  token        = "${chomp(file(var.scaleway_token_file))}"
  region = "${var.scaleway_region}"
}

module "scaleway" {
  source = "./scaleway"
  machine_name = "prod.hkjn.me"
  image = "ee0d3a38-1e8a-4407-bc02-d35dd588efa2"
}

module "ln_scaleway" {
  source = "./scaleway"
  machine_name = "ln.hkjn.me"
}
