variable "gcloud_credentials" {
  default = ".gcp/tf-dns-editor.json"
}

variable "hkjnprod_enabled" {
  default = true
}

variable "scaleway_organization_file" {
  default = ".scaleway/scaleway0_organization"
}

variable "scaleway_token_file" {
  default = ".scaleway/scaleway0_token"
}

variable "scaleway_region" {
  default = "par1"
}

variable "digitalocean_token_file" {
  default = ".digitalocean/digitalocean0_token"
}

variable "digitalocean_image" {
  # Images can be fetched with:
  # curl -X GET --silent "https://api.digitalocean.com/v2/images?per_page=999" -H "Authorization: Bearer $(cat .digitalocean/digitalocean0_token)" |jq '.'
  default = "coreos-stable"
}

variable "admin2_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDb7dnRu7p1esBjqNbt3pHelNtu+yJqkFUt0behdnOk+gnU/EjMraVwdmcr0EZfXht6mettQAGil3TsT0NznUlpG6xAKigzLxhBt43Ul9Q9A6XDF04+7vRH/V/TSV6P5nYREgl9NthAJwpXGTijnapZcz3fhmQMuH3Xg1+H+2Lh8MdaduN7cVXDlSht1r8sJSvQKEl7jr7mp4soegw9x/9regsRqD7cxGfcXZ22PuedD4M0HrT4B4Sva5vQKq6WcqujyLLNStcYRBEmUQgIkG528XkGvxqaOOGMIhwPwln+TPat1nXMjGk9pGMNqawRYNuL/dhn18JOpsVecPfrMNr9H9Kcjoi1BdiEjs7rtEyy2i7knuosR+7pS8gNuyQWZvLzsQ80t1yyxGJ0pvxki/Zijp5vzwluUzqoBebzD0YvSb9u1oryjRMUtum6jTay+GbzKYyMmggAhYH8Eog4bllUhjrdy0swwka2TXnot/RG+Uueos5caD+EfKICem/P16zHwHWy+9cm3ZUOXPVG6yjPs1ZjNHy2YJ5OkbRbNHRpqaNNA8CYJLNpGG5J46aVRQMbi9dD94cK7B2fA+HSABu5ttTPlciBD6DrS6HXeZJhbg6gHpT5mQYjxMEaBOzoKzB0oILkUPrCMPyXQtvQtheASh7fPePezkwNM37Y9OHYAQ== admin2_id_rsa@ipad"
}

variable "digitalocean1_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCtb2pbWvs+LZapn7yTtKD0NrTBmgTXoDgERZAcJ47ziGawmaBtQ6UgZLppJ10atm/SddQohYwRNnZp1XukbQf83g8XgJ+W91WIHZdggmjnVIGIan+TbNsNVPAtpRXlEIiYSgpy8jHci1q0u6jat2/FGi0x3AUqGUAaUJxwvUEDnkf/g25lq/hyOZi0yoCFjytQC3/TlgFaCW4T/8RRPY4aHnEtY/D0GE4UPBsDBK+wT2/cFrcxAVSNLHW5i44ChzXbCtwTyPv9+FayngjQyPtze84KPZa7gv7XITcAmOnfRRNCN5UlNCJEwMXg2jZtext4OuUPcZo1z5/D6GaP3n+3WALIk87h8+a4aJ0hJ3RwSNEZaKicgTsziNfmuXaztC/DZAG/fpbdG0O+VULwHwTIFJVsXB9yP2bXqmJOTj9/T0NimN8XPzBa+ixo7BPebuMCwyIS31zhdzxoidi3tt/bHCozY0Aoh/sGelku0xIsT8Io3WrMX5Cqgbz3EPnjeMYP9kbG695xwRfRdS0t+qjxezG73wxs4FPwPY4e9T19G8v91XzgbKK9c33S4DyoL4Zf4nGW+i0dqkBsLgegMLXAXsNiwA1hcHDLXm7hO+viwC50RSwHUeCvwNPDzX18QSmw+mt+OOoMClnWdM6IuOzO3q4TwkwMkpqY8ncrmg+8XQ=="
}

variable "aruna_pubkey_arunallave" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC0qu+LwXD4QG1kQv4F8wjY1bIDR8J72kX/bt8kP/qdYJ0Y0dzSjwGUG79BVBx/StkDRYr+eiDk/fa/mL5ZUCtsraRiOyPoCqvfkIBvHRAjiKynwa/gTCHw1gC/fT2HgTq996VNoqnZDLV4Icc9t7BPMecaNVquBc3+eW6N8fCSV1Oj8/+NlJyV2Efom2ciG/GURknYO5DdhcH+NqmvjhMlBy/XLkfCSGiwNXm4XRHDAe5ZKELTDHH4sFXCPUjSbv7tB75iHDY71lBU+yNb7i3uOO/ODluCxeXranWbVyE76IpA8XoAPdsq0xcfZKLt3ZEGd1NU+v8ZWMXrTUJxX147UGoXsh76pPnY3pgrDxY75K8zgetS2y/TsNcnZP4KvWZK/6hVGasyhmoQxYgr4yj+4L8J9izzGSYYmm+vXL8Q+1ejgWG34KzRJjK8kQvdAKiSu+nnnPSoizsv6WEowAkBQp4pfqjHf5LIFHrfvVoU7HFLKp0wWSYeF/3chl65rhu8y67TkWHBq3mVqCEKqnciXgbYu/qMlfWpv3uHqi1Qhodpm3PPM3O1rk/+WKKrfi1gwcr6F+tz6aq6y1qnTx4YeW0Bqd1uf7J5vEe0HKjjL3Dt9BsW9f5tS3FAic+AtUlPO2vADYYonYc9FiLlux/UCDS7oIBUjsfczsM8ML0M9w== aruna@ipad"
}

variable "admin1_blockpress_me_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC5P2mi5BvGKR4qv6+BFt9N1Dnom2NUo/uIXpZHAmeHgP7RC4gL2LXYGpIFphRuEvDzPaCuIK98Cqu1AUKCwoIx3MzDeK55u2dbYlvs2VJC0lUPzj5nm6OHWKTrfsjtkBCqv2DjWKrCNFNVvBV3hLPKThS5kmNh8RO8zbxk6LsBw4tJYc7ZYgL/0USbqq2Ud/Zr2bjspaGLim1mshQo+tOhu8sQa5K7JpZYKSeAajnAj9U/TeWwGIK1wITrzCDpYroKNT3vMTTxxhlP1aARAJTPt+E8wDt377GFJ1u+qte2LTLyC84U6pw1kFCz1r/NRuqlXcs98pu/L+rIkTEcraHDdD3KITyW/4KkUOkt1SEV14TZH+JZ7jQ1vL6ztFeaduTL95H7YrprJWds2/SZOJPrxSusYjibjCSddybJWXmPFleWSNqs7sAQ+ljsStqm1BIinR7L++FRacws/D3Si4cgEkP9TZAaLhHjW+MF1e0ueB/KVgvEZ45uCJQdefxRDjXs6VYlKUVg59NW0fzsNi5ro3NEEyg6Sgn/w34JCyhgz41h5ncdZZKF/Ti6n40onwuuiV73RXec4ViVrW6vJOUGwZFVPX11ra3KcSgyNHGPq6a0lNmjfmBkzku4gaq8ZKluIam4NyIvGkcacA3XAcUfLD1N9RwBGAO8TVAff9C75w== zero@admin1"
}

variable "admin2_builder_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDDt1kEy5DIwmIWFLlziyDiB0LIU3cyXTVGsPk07Y1oXmHVr84RX27geHw1VtnFAj0PJH2z+yk08Xd+KQxxwb8Ijv3rv+RG1hxgSioLc5OuqblD63xTBoA9MP4pXhaVnItbuvfcivzJXdP5FCVEn+Wfuw0g82Rr516QpsL22D4TRdPcrpzR3WbO48bfZNKKq8uOUcahEwcFbD0T9W6brCXf9vKOZTzyXpapfJ/pBPE0PnHZ+c1kEs0dMqYk4iZsY8RkWq0uNelsI8WcVwRsykvzT7CBm0qxr8O+q+j59L0oDJ3kZ3L33OkOqBD1i++jYQhQ6hTtgq4UWdIxolN0DhidETXYEjyFALLL/RD2ba8pIaBRm/pmrsyPhQJFG8oIo5FIoxfyM0Kf4IOA42ofzx2ukgJCrwbNfSC6oRBfuCw8O/YBuj1pi3/X3qArmeGaxOvrKUUYBC7cApqJ4czFyg6xOpNxv2pFNAcTt1C7st+7NsExNZ32VLSFnoPi/DgcRiu47I+ianLE0rerceCH2BJxlKc+Fb1sz9uXympvCmxMbd+3bJa1pbtX01E6sjI7h80NNm573UfHa51hUJtopnXrGzcOkaDbaZRxZQbvVtEcRoG/jB2KiiIenM6CCJl4XdYnAIySIF5++yDiyZX9BnHuGKyH+h6Wbjk3Cfvm3ow0qQ== zero@admin2"
}

variable "admin2_decenter_world_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCsgXS6D4sn75MmcsG0chzXxwNnyJQJm6IfKSQSOCVZt6X7Js0k+D/J8vM2tDQzlQYNGWy1Qk+JQk0B8XZrImpiYYBtVUbRmaqD6mHNrT6vdovdh8piII0dNJY0MHIgckyK806EZurbcbApyVGlK5xZTILuEM5h5Iduhhu5HnHL0zzvjjAVk1C6CtbwKUmnfxZaHztN0ayJf6PYmQ+yJBQha+6+BiJm1NfUtMtn8K+fX7fIHYTcTUMNJOvGRRZXZ5LD+4s4XiIFKO2BH2XjZVhUawVpE6Z2okbuotrcKqcS4DipbobXnB7ujTTE4fJwzWuGH8lAH7/XSUIaV8NKLqH6osm3ee4OBhEa3MInQfAAYB6IkroxSKbuYvQiroG2R+vypSWgNCEX6bLzInR8Bddjn/saq7nzYtm0fybPlsHqrsb8y14yD0TAPb5swQKnsu561oExU5vEcIEum/uVJTIurT2qbOo6U+OT8tN8mR0y4aEdYDEwBqSa8BPKCgG+8KPri9nKotDzXeD1lxdXL/0ZeuREIwwAeFDkjhID1KUAXRdCysC6IsQ4W/nSeZkhldqcqpEsLhU/n9qYDEoJ/jW576q32zo9Jr55Gvgo3wMUfiyuPJ5Q8oi2Iv9Anl/1iMcr83rPsnh21CzVBnPoTmSHAvYHv0YzQVZ5ZcTvMmAzmw== zero@admin2"
}

variable "admin2_guac_pubkey" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDbEY3VQRAAjlnI31tvGFgnJaYNdipRH6LFlBir646WEDp+xZWI2BPPdy25qLOovGko5fVegLC7DISzTz1RMmzyVgd8IMpK5wKcHOK/FR2E3nRqVg2PMsFkHrboWuFDqFyLBvoq5wq+3lN86tDIhc+QZOT077cjBfeRrm2CDqEh41BqcIwXF+zAYI14NjreK/A2ug03vBWv1akJ3lonuED2+VCKn1afyf9w8YUlQZ5Hc+X6yGPPgZkCVaSiZWMwFBHI87h8HMjACUSTeRbkdugmrwtcmmPTLMPJgvW7pdVBbYOhLtaQLYqz5ar2p8mWNQMdF/WSObu5iJnnQ/9bt6QBHNOcB/bdtb6orT67tIV42rg4Pm1YpqfDj7hzJiMlEgU4LBpxGPuYD/BhZdSMZB6tnl43LziqzassbKtq1kMQTxX4Y+r0XTVaUDFg+mpa+8A8MlMJLNjPvEaanxgGOpW/mOtrv1Ta5zXRuwvBxMa5RI3nvoseqLSwb6psYhoZyhLDxjRAVf5y4z8JJX1+q2VCXWREcsuPnFCnpdnaBJzaUp0L9Cz+PXCkN5VlndJpdTH+8fB8ijypiCWU6hs3gzRy/Lu7d1FuU6MRGbEfyvzuUv9EwDi4cFhx0bsgx/Ymzx1wVKxFKhCWfxPXEc7gA4JIrKvHt8IAjch8Uu0x6wMkaQ== zero@admin2"
}

variable "gcloud_project" {
  default = "henrik-jonsson"
}

variable "gcloud_region" {
  default = "europe-west1"
}

#
# The latest image can be found with:
# $ gcloud compute images list | grep coreos-alpha
#
variable "coreos_alpha_image" {
  default = "coreos-alpha-1632-0-0-v20171220"
}

#
# The latest image can be found with:
# $ gcloud compute images list | grep coreos-beta
#
variable "coreos_beta_image" {
  default = "coreos-beta-1492-6-0-v20170906"
}

#
# The latest image can be found with:
# $ gcloud compute images list | grep ubuntu
#
variable "ubuntu_image" {
  default = "ubuntu-1710-artful-v20171122"
}

variable "admin1_ip" {
  default = "51.15.200.169"
}

variable "exocore_ip" {
  default = "159.100.250.108"
}

variable "hkjnweb_ip" {
  default = "163.172.173.208"
}

variable "mon_ip" {
  default = "163.172.184.153"
}

variable "vpn_ip" {
  default = "163.172.184.153"
}

variable "cities_ip" {
  default = "163.172.184.153"
}

variable "builder_enabled" {
  description = "Whether builder node is enabled."
  default = false
}

variable "dropcore_enabled" {
  description = "Whether dropcore node is enabled."
  default = false
}

variable "elentari_world_enabled" {
  description = "Whether elentari.world node is enabled."
  default = false
}

variable "guac_enabled" {
  description = "Whether guac.hkjn.me node is enabled."
  default = false
}

variable "blockpress_me_enabled" {
  description = "Whether blockpress.me node is enabled."
  default = true
}

variable "version" {}
