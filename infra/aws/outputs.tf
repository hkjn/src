output "lab_eip" {
	value = "${aws_eip.eip.public_ip}"
}
