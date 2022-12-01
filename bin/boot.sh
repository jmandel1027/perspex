#!/usr/bin/env bash

set -e

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
  
  bin/helm dependency update infrastructure/charts/perspex

  echo ""
}

function activate_cluster() {
  
  ping_docker

  upgrade_helm_deps

  echo ""

  echo "Verifying k8s cluster is provisioned"
  if ! bin/k3d cluster list | grep -o -q "${cluster_name}" &> /dev/null ; then 
    bin/k3d cluster create -c infrastructure/tilt/k3d-config.yaml
  fi

  echo ""

  echo "Spinning up tilt"
  bin/tilt up
  echo ""
}

function deactivate_cluster() {
  echo "Tearing down k8s cluster"
  if ! bin/k3d cluster list | grep -o -q "${cluster_name}" &> /dev/null ; then 
    bin/k3d cluster delete "${cluster_name}"
  fi

  echo ""
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
    if [ "$(bin/kubectl get ns "${NAMESPACE}" -o json | jq .status.phase -r)" == "Active" ]; then
        break
    fi
    echo "Namespace ${NAMESPACE}" not ready
    bin/kubectl create ns "${NAMESPACE}"
  done

  echo "Namespace ${NAMESPACE} ready - creating local secrets"

  bin/kubectl create secret generic perspex-local-secrets -n "${NAMESPACE}" --save-config --dry-run=client --from-env-file=.env -o yaml | kubectl apply -f -
}

function teardown() {
  echo "Tearing down ${cluster_name}"
  bin/k3d cluster delete "${cluster_name}"
  echo ""

  exit 0
}

function display_help() {
    echo "This script will run a database migration validation service"
    echo ""
    echo "Parameters:"
    echo "    -o | --on               Activate the dev environment."
    echo "    -f | --off              Deactivate the dev environment with out teardown."
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