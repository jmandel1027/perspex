#!/usr/bin/env bash

set -e


if k3d cluster list | grep -o -q perspex-local &> /dev/null ; then 
  k3d cluster create  -c infrastructure/tilt/k3d-config.yaml
fi