# infra

This directory holds infrastructure configs and tools for the hkjn.me site and other projects.

## Dockerized Terraform

A dockerized Terraform alias `tf` can be added to our shell by:

```
source tf_dockerized.sh
```

After this, we can use `tf plan`, `tf apply` and other commands.

## Dockerized gcloud tools

An alias for working with a dockerized set of `gcloud` tools can be
added by the command:

```
source gcloud_dockerized.sh
```

After this, the `gcd` command enters an interactive container which
can work towards the GCE project.

## Ignition format of  `user-data`

The `user-data` field for CoreOS machines is in Ignition format:

* https://coreos.com/os/docs/latest/booting-on-google-compute-engine.html

The `generate_ignite_configs.go` tool generates Ignition `.json` configs for
nodes that make use of Ignition to know what tasks should be done on first boot.

## Tests

The `run_tests` script runs all relevant tests. It can be added to `git`
pre-push and pre-commit hooks by doing:

```
cd .git/hooks
ln -s ../../run_tests pre-commit
ln -s ../../run_tests pre-push
