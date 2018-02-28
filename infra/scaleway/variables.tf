variable "enabled" {
  default = true
}

variable "image" {
  # Fetched with:
  # curl -H "X-Auth-Token: $(cat .scaleway/scaleway0_token)" -H 'Content-Type: application/json' 'https://cp-par1.scaleway.com/images/?page=1&per_page=100' > scaleway_images.json
  # >>> import json; d=json.loads(open('images.json').read())
  # >>> [x for x in d['images'] if x['arch'] == 'arm' and 'Docker' in x['name']]
  # TODO: Find reason why machine doesn't allow ssh with following image:
  default = "2abb57bb-427d-4317-85f0-799f247a2224"
}

variable "machine_name" {
}
