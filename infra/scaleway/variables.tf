variable "enabled" {
  default = true
}

variable "image" {
  # Fetched with:
  # curl -H "X-Auth-Token: $(cat .scaleway/scaleway0_key)" -H 'Content-Type: application/json' 'https://cp-par1.scaleway.com/images/?page=1&per_page=100' > scaleway_images.json
  default = "ee0d3a38-1e8a-4407-bc02-d35dd588efa2"
}

