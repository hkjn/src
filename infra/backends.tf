terraform {
  backend "gcs" {
    bucket  = "hkjn-terraform-state"
    prefix    = "hkjninfra/prod/terraform.tfstate"
  }
}
