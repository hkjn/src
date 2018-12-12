variable "image" {
  # Fetched with:
  # curl -H "X-Auth-Token: $(cat .scaleway/scaleway0_token)" -H 'Content-Type: application/json' 'https://cp-par1.scaleway.com/images/?page=1&per_page=100' > scaleway_images.json
  # >>> import json; d=json.loads(open('images.json').read())
  # >>> [x for x in d['images'] if x['arch'] == 'arm' and 'Docker' in x['name']]
  # name: 'Ubuntu Xenial (16.04 latest)'
  default = "3a1b0dd8-92e1-4ba2-aece-eea8e9d07e32"
}

variable "machine_name" {}
