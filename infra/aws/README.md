# aws

This subdirectory holds terraform plans for infrastructure on AWS.

## setup

After a clean checkout, AWS API keys need to be specified by creating `keys.tf`:

```
variables "aws_access_key" {
	default = "<ACTUAL AWS ACCESS KEY GOES HERE>"
}
variables "aws_secret_key" {
	default = "<ACTUAL AWS SECRET KEY GOES HERE>"
}
```

It might also be necessary to change the ownership of the directory to match
the user inside the `hkjn/terraform` image:

```
$ sudo chown -R 1000:1000 .
```
