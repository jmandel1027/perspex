# Perspex

## Getting started

If this is your first time running the project you'll need to install `helm`, `k3d` and `tilt` to set up local development environment

This command will spin up the entire stack, verifying that docker is active and install dependencies like k3d, helm and tilt if their not already present

```sh
bin/boot.sh -o
```



```sh
    cd perspex
    tilt up
```