# building with docker instead of podman

```
BUILDER=docker make
```

# time drift issues on macos

needed to run following to do local build (without ability to pull
from docker registry, due to podman VM time drift issue on macos):

```
podman tag hkjn/alpine:0.2.0-arm64 hkjn/alpine:0.2.0
```

workaround until https://github.com/containers/podman/issues/11541
is fixed is:

```
podman machine ssh --username root date --set $(date -Iseconds)
```
