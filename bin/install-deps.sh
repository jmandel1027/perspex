#!/usr/bin/env bash

set -e

dir=$(pwd)

function install_deps() {
    
    # check if helm is installed
    if ! command -v helm &> /dev/null ; then
      brew install helm

      cd infrastructure/charts/perspex
      helm dep update

      cd "${dir}"
    fi

    # check if k3d is installed
    if ! command -v k3d &> /dev/null ; then
      # this can be found via https://k3d.io/v5.2.2/
      curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
      
    fi

    if ! command -v tilt &> /dev/null ; then

      # this can be found via https://docs.tilt.dev/install.html
      curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
    fi

}


install_deps