provider "aws" {
	access_key = "${var.aws_access_key}"
	secret_key = "${var.aws_secret_key}"
	region     = "eu-west-1"
}

data "aws_ami" "ubuntu" {
	most_recent = true

	filter {
		name   = "name"
		values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
	}

	filter {
		name   = "virtualization-type"
		values = ["hvm"]
	}

	owners = ["099720109477"] # Canonical
}

resource "aws_security_group" "allow_ssh" {
	name        = "allow_ssh"
	description = "Allow inbound ssh traffic"
	vpc_id      = "${var.vpc_id}"
}

resource "aws_security_group" "allow_https" {
	name        = "allow_https"
	description = "Allow inbound https traffic"
	vpc_id      = "${var.vpc_id}"
}

resource "aws_security_group" "allow_outbound" {
	name        = "allow_outbound"
	description = "Allow all outbound traffic"
	vpc_id      = "${var.vpc_id}"
}

resource "aws_security_group_rule" "allow_ssh" {
	type            = "ingress"
	from_port       = 22
	to_port         = 22
	protocol        = "tcp"
	cidr_blocks = [
		"35.198.109.207/32",
		"77.56.54.251/32",
		"178.197.228.188/32",
	]
	security_group_id = "${aws_security_group.allow_ssh.id}"
}

resource "aws_security_group_rule" "allow_dev" {
	type            = "ingress"
	from_port       = 8888
	to_port         = 8888
	protocol        = "tcp"
	cidr_blocks = [
		"77.56.54.251/32",
	]
	security_group_id = "${aws_security_group.allow_ssh.id}"
}

resource "aws_security_group_rule" "allow_http" {
	type            = "ingress"
	from_port       = 80
	to_port         = 80
	protocol        = "tcp"
	cidr_blocks = ["0.0.0.0/0"]
	security_group_id = "${aws_security_group.allow_https.id}"
}

resource "aws_security_group_rule" "allow_https" {
	type            = "ingress"
	from_port       = 443
	to_port         = 443
	protocol        = "tcp"
	cidr_blocks = ["0.0.0.0/0"]
	security_group_id = "${aws_security_group.allow_https.id}"
}

resource "aws_security_group_rule" "allow_outbound" {
	type            = "egress"
	from_port       = 0
	to_port         = 0
	protocol        = "-1"
	cidr_blocks = ["0.0.0.0/0"]
	security_group_id = "${aws_security_group.allow_outbound.id}"
}

resource "aws_eip" "eip" {
	instance = "${aws_instance.lab.id}"
	vpc      = true
}

resource "aws_instance" "lab" {
	ami             = "${data.aws_ami.ubuntu.id}"
	instance_type   = "t2.medium"
	key_name        = "zaws_key1"
	security_groups = [
		"${aws_security_group.allow_ssh.name}",
		"${aws_security_group.allow_https.name}",
		"${aws_security_group.allow_outbound.name}"
	]
	root_block_device = {
		volume_size = 25
	}

	tags {
		Name = "lab0"
	}
}
