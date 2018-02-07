#
# *.decenter.world
#
output "decenter_world_addr" {
  value = "${google_dns_record_set.decenter_world.name}"
}

#
# *.blockpress.me
#
output "blockpress_me_addr" {
  value = "${google_dns_record_set.blockpress_me.name}"
}

output "blockpress_me_ip" {
  value = "${google_dns_record_set.blockpress_me.rrdatas[0]}"
}

#
# *.elentari.world
#
output "elentari_world_addr" {
  value = "${google_dns_record_set.elentari_world_web.name}"
}

output "elentari_world_ip" {
  value = "${google_dns_record_set.elentari_world_web.rrdatas[0]}"
}

#
# *.hkjn.me
#
output "hkjn_addr_builder" {
  value = "${google_dns_record_set.hkjn_builder.name}"
}

output "hkjn_addr_dropcore" {
  value = "${google_dns_record_set.hkjn_dropcore.name}"
}

output "hkjn_addr_exocore" {
  value = "${google_dns_record_set.hkjn_exocore.name}"
}

output "hkjn_addr_admin1" {
  value = "${google_dns_record_set.hkjn_admin1.name}"
}

output "hkjn_addr_admin2" {
  value = "${google_dns_record_set.hkjn_admin2.name}"
}

output "hkjn_addr_cities" {
  value = "${google_dns_record_set.hkjn_cities.name}"
}

output "hkjn_addr_guac" {
  value = "${google_dns_record_set.hkjn_guac.name}"
}

output "hkjn_addr_mon" {
  value = "${google_dns_record_set.hkjn_mon.name}"
}

output "hkjn_addr_vpn" {
  value = "${google_dns_record_set.hkjn_vpn.name}"
}

#
# IP addresses of *.hkjn.me
#

output "builder_ip" {
  value = "${google_compute_instance.builder.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "dropcore_ip" {
  value = "${digitalocean_droplet.dropcore.ipv4_address}"
}

output "hkjn_prod_ip" {
  value = "${module.scaleway.public_ip}"
}

output "hkjn_ip_gz0" {
  value = "${google_dns_record_set.hkjn_gz0.rrdatas[0]}"
}

output "hkjn_ip_zs10" {
  value = "${google_dns_record_set.hkjn_zs10.rrdatas[0]}"
}

output "hkjn_ignite_json_decenter_world" {
  value = "${google_compute_instance.decenter-world.metadata.user-data}"
}

output "hkjn_ip_guac" {
  value = "${google_compute_instance.guac.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "hkjn_ip_decenter_world" {
  value = "${google_compute_instance.decenter-world.network_interface.0.access_config.0.assigned_nat_ip}"
}

output "hkjn_ip_web" {
  value = "${var.hkjnweb_ip}"
}
