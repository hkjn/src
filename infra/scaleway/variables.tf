variable "image" {
  # Fetched with:
  # curl -H "X-Auth-Token: $(cat .scaleway/scaleway0_token)" -H 'Content-Type: application/json' 'https://cp-par1.scaleway.com/images/?page=1&per_page=100' > scaleway_images.json
  # >>> import json; d=json.loads(open('scaleway_images.json').read())
  # >>> pprint.pprint([x for x in d['images'] if x['arch'] == 'arm' and 'Ubuntu Xenial' in x['name']])
  default = "8c76c2b9-926e-44c5-b250-6f6480b5d313"
}

