#!/usr/bin/env bash

set -e

dir=$(pwd)
cluster_name="perspex-local"


function ping_docker() {
  echo "Verifying that docker daemon is active"
  if (! docker stats --no-stream ); then
    # On Mac OS this would be the terminal command to launch Docker
     open "/Applications/Rancher Desktop.app" || open "/Applications/Docker.app"
    # Wait until Docker daemon is running and has completed initialisation
    while (! docker stats --no-stream ); do
      # Docker takes a few seconds to initialize
      echo "Waiting for Docker to launch..."
      sleep 1
    done
  fi
  echo ""
}


function upgrade_helm_deps() {
  echo "Upgrading chart deps"
  
  cd infrastructure/charts/perspex
  helm dep update

  cd "${dir}"
  echo ""
}

function activate_cluster() {
  
  ping_docker

  upgrade_helm_deps

  echo "Verifying deps are installed"
  if ! command -v tilt &> /dev/null ; then
    install_deps
  fi

  echo ""

  echo "Verifying k8s cluster is provisioned"
  if ! k3d cluster list | grep -o -q "${cluster_name}" &> /dev/null ; then 
    k3d cluster create -c infrastructure/tilt/k3d-config.yaml
  fi

  echo ""

  echo "Spinning up tilt"
  tilt up
  echo ""
}

function deactivate_cluster() {
  echo "Tearing down k8s cluster"
  if ! k3d cluster list | grep -o -q "${cluster_name}" &> /dev/null ; then 
    k3d cluster delete "${cluster_name}"
  fi

  echo ""
}


function install_deps() {
  if ! command -v kubectl &> /dev/null ; then
    echo "Installing kubectl"
    brew install kubectl
  fi

  # check if helm is installed
  if ! command -v helm &> /dev/null ; then
    echo "Installing helm"
    brew install helm
    echo ""

    echo "Installing chart deps"
    cd infrastructure/charts/perspex
    helm dep update

    cd "${dir}"
    echo ""
  fi

  # check if k3d is installed
  if ! command -v k3d &> /dev/null ; then
    echo "Installing k3d"
    # this can be found via https://k3d.io/v5.2.2/
    curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash    
    echo ""
  fi

  if ! command -v tilt &> /dev/null ; then
    echo "Installing tilt"
    # this can be found via https://docs.tilt.dev/install.html
    curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
    echo ""
  fi

  if ! command -v migrate &> /dev/null ; then
    echo "Installing migration util"
    brew install golang-migrate
    echo ""
  fi

  exit 0
}

create_local_sectets() {
  # CHECK IF .ENV EXISTS
  # IF NOT CREATE .ENV BY
    # SEED FROM HERE DIRECTLY?
    # SEED BY PARSING VALUES IN TILT?
  if [[ ! -f ".env" ]]; then
    echo ""
  fi
  echo ""
}

generate_local_kube_secrets() {
  NAMESPACE=${1}
  MAX=10
  COUNTER=0

  while [ $COUNTER -le $MAX ]; do
    if [ "$(kubectl get ns "${NAMESPACE}" -o json | jq .status.phase -r)" == "Active" ]; then
        break
    fi
    echo "Namespace ${NAMESPACE}" not ready
    kubectl create ns "${NAMESPACE}"
  done

  echo "Namespace ${NAMESPACE} ready - creating local secrets"

  kubectl create secret generic perspex-local-secrets -n "${NAMESPACE}" --save-config --dry-run=client --from-env-file=.env -o yaml | kubectl apply -f -
}

function teardown() {
  echo "Tearing down ${cluster_name}"
  k3d cluster delete "${cluster_name}"
  echo ""

  exit 0
}

function display_help() {
    echo "This script will run a database migration validation service"
    echo ""
    echo "Parameters:"
    echo "    -o | --on               Activate the dev environment."
    echo "    -f | --off              Deactivate the dev environment with out teardown."
    echo "    -i | --install          Install dependencies"
    echo "    -t | --teardown         Teardown development infrastructure."
    echo "    -h                      Display this help message."
    echo "    -b                      Specifies the branch to validate migration against main"
    echo "Usage:"
    echo "   bin/boot.sh -o"
    echo "   bin/boot.sh -f"
    echo "   bin/boot.sh -i"
    echo "   bin/boot.sh -t"
    echo "   bin/boot.sh -h"
}

main() {
  while [[ "$#" -gt 0 ]]; do 
    case $1 in
      -d | --destroy )
        teardown
        shift
        ;;
      -o | --on )
        activate_cluster
        shift
        ;;
      -i | --install )
        install_deps
        shift
        ;;
      -h | --help )
        display_help
        shift
        ;;
      \? )
        display_help
        exit 1
        ;;
    esac;
    shift;
  done
}

main "${@}"