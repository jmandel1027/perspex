#!/usr/bin/env bash

set -e

connect_to_postgres() {
  if [[ ${environment} == "development" ]]; then
    echo "Info: accessing the dev cluster's DB will prompt a change in kubectl's ctx"
    kubectl config use-context development-eks-cluster

    cmd="psql postgres://postgres:$POSTGRES_PASSWORD@postgresql:5432/perspex"
    kubectl -n portal exec -it pod/portal-postgresql-0 -- sh -c "$cmd"
  elif [[ ${environment} == "tilt" ]]; then
    echo "Info: accessing the dev cluster's DB will prompt a change in kubectl's ctx"
    kubectl config use-context k3d-perspex-local

    cmd='psql "postgres://perspex:pass@postgresql:5432/perspex"'
    kubectl -n perspex exec -it pod/postgresql-0 -- sh -c "$cmd"
  fi

  exit 0
}

usage() {
    echo "Usage: $0 [option...]" >&2
    echo
    echo "   -local, --localhost        Connect to postgres locally"
    echo "   -dev, --development        Connect to postgres in development"
    echo "   -stag, --staging           Connect to postgres in staging"
    echo "   -prod, --production        Connect to postgres in production"
    echo
  exit 0
}

main() {
  while [[ "$#" -gt 0 ]]; do
    case $1 in
      -local|--localhost)
        environment="localhost"
        ;;
      -tilt|--tilt)
        environment="tilt"
        connect_to_postgres
        ;;
      -dev|--development)
        environment="development"
        ;;
      -stag|--staging)
        environment="staging"
        ;;
      -prod|--production)
        environment="production"
        ;;
      -h|--help)
        usage
        ;;
      *)
        echo "Unknown parameter passed: $1";
        exit 1
      ;;
    esac;
    shift;
  done
}

main "${@}"