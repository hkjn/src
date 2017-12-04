terraform {
  backend "gcs" {
    bucket  = "hkjn-terraform-state"
    path    = "hkjninfra/prod/terraform.tfstate"
    project = "henrik-jonsson"
  }
}
