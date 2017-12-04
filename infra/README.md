# hkjninfra

Repo hkjninfra holds infrastructure configs for hkjn.me and other projects.

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

The `ignite.py` script generates the Ignition `.json` configs, which
tells the instances what tasks should be done on first boot:

```
python ignite.py
```

## Tests

The `run_tests` script runs all relevant tests. It can be added to `git`
pre-push and pre-commit hooks by doing:

```
cd .git/hooks
ln -s ../../run_tests pre-commit
ln -s ../../run_tests pre-push
```
