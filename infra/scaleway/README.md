# scaleway

This subdirectory holds terraform plans for infrastructure on scaleway.com.

## setup

After a clean checkout, scaleway.com API keys need to be specified by creating `keys.tf`:

```
variables "scaleway_organization" {
	default = "<ACTUAL SCALEWAY ORGANIZATION GOES HERE>"
}
variables "scaleway_token" {
	default = "<ACTUAL SCALEWAY TOKEN GOES HERE>"
}
```

The credentials for the backend storage for the Terraform state file also needs to
be set in the `.backend_credentials` file:
```
[default]
aws_access_key_id = <ACTUAL AWS ACCESS KEY GOES HERE>
aws_secret_access_key = <ACTUAL AWS SECRET KEY GOES HERE>
```

It might also be necessary to change the ownership of the directory to match
the user/group id used inside the `hkjn/terraform` image:

```
$ sudo chown -R 1000:1000 .
```
