# Perspex
[![Go Reference](https://pkg.go.dev/badge/github.com/jmandel1027/perspex.svg)](https://pkg.go.dev/github.com/jmandel1027/perspex)

![Alt text](.github/media/perspex.jpg?raw=true "perspex-icon")


## Table of Contents

- [About](#about)
  * [Core Concepts](#core-concepts)
  * [Architecture](#architecture)
- [Getting Started](#getting-started)
  * [Manual Commands](#manual-commands)
- [Project Structure](#project-structure)
- [Adding Dependencies](#adding-dependencies)


## About

Perspex is something like a sourdough starter, it's aimed to be a completely portable development envirionment (atleast for unix-like platforms, sorry windows) and toolchain for gRPC, GraphQL, Tilt, K8s, and more.
Utilizing [hermit](https://cashapp.github.io/hermit/) we can ensure that anyone who pulls down this codebase can immediately spin up the entire stack without worrying about dependency conflicts. Better still we can consistently ensure that the behavior is 1:1 throughout every cycle of development since the CI/CD processes utilize the exact same deps as we utilize locally.

### Core Concepts
  - Hermeticity: System deps should be managed by hermit to ensure no barries to boot
  - Schema Driven: Design your API how you want to express your API and DB operations, not the other way around.
  - Consistent Topologies: a common Kubes based deployment interface to unify local dev, CI, testing and prod.

### Architecture

![Alt text](.github/media/perspex-arch.png?raw=true "perspex-icon")

This is a high level overview of the current system design. We have inbound traffic coming from either a React frontend or via autogenerated gRPC Client SDKs. We 

## Getting started

This command will spin up the entire stack, verifying that docker is active. Since this project utilizes [Hermit](https://cashapp.github.io/hermit/) for hermetic dependencies no external dependencies other than docker are required. 

```sh
bin/boot.sh -o
```

### Manual Commands

For manually controlling the stack the following commands are here to help get things rolling.

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


## Project Structure

```
├── .github
├── README.md
├── bin
|  ├── activate-hermit
|  └── etc ...
├── schema
|  ├── graphql
|  ├── proto
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

## Adding Dependencies

If you need to add deps to hermit please perform the following.

```sh

cd perspex

. bin/activate-hermit

hermit search terraform 

hermit install terraform-x.x.x

deactivate-hermit
```

If you need to add a dependency to the gopath, please perform the following

```sh
cd perspex

bin/go install google.golang.org/protobuf/cmd/protoc-gen-go@latest 

```
