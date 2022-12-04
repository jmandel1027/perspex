# Perspex

![Alt text](.github/media/perspex.jpg?raw=true "perspex-icon")

## What is this?

I'm not entirely sure, just experimenting with various toolings and infrastructure. I'm mostly just having fun and thowing paint at the wall.

## Getting started

This command will spin up the entire stack, verifying that docker is active. Since this project utilizes [Hermit](https://cashapp.github.io/hermit/) for hermetic dependencies no external dependencies other than docker are required.

```sh
bin/boot.sh -o
```

```sh
# this installs third party helm deps
bin/helm dependency update infrastructure/charts/perspex

# provisioons the k3d cluster and registry
bin/k3d cluster create -c infrastructure/tilt/k3d-config.yaml

# spin up tilt
bin/tilt up
```

To teardown the entire stack, please run the following:
```
bin/boot.sh -d
```
##

Project Structure

```
├── .github
├── README.md
├── bin
|  ├── activate-hermit
|  └── etc ...
├── schema
|  ├── graphql
|  └── etc ...
├── services
|  ├── backend
|  └── etc ...
├── infrastructure
|  ├── charts
|  |  ├── perspex
|  |  └── etc ...
|  ├── terraform
|  |  └── modules
|  ├── terragrunt
|  └── etc ...
└── Tiltfile
```
