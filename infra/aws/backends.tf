terraform {
  backend "s3" {
    bucket = "zterraform-state"
    key    = "hkjn.me/src/infra/aws/terraform.tfstate"
    region = "eu-west-1"
    shared_credentials_file = ".backend_credentials"
  }
}