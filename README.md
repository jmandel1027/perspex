# Perspex

## Getting started

If this is your first time running the project you'll need to install `helm`, `k3d` and `tilt` to set up local development environment

This command will spin up the entire stack, verifying that docker is active and install dependencies like k3d, helm and tilt if their not already present

```sh
bin/boot.sh -o
```
If you would like to operate the stack manually and have tilt, k3d, and helm installed, please perform the following 

```sh

dir=$(pwd)

cd infrastructure/charts/perspex

# this installs third party helm deps
helm dep update

cd "${dir}"

# provisioons the k3d cluster and registry
k3d cluster create -c infrastructure/tilt/k3d-config.yaml

# spin up tilt
tilt up
```

To teardown the entire stack, please run the following:
```
bin/boot.sh -d
```