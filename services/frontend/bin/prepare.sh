#!/bin/bash

set -e

build_proto() {
  cd ../../schemas/proto

  ../../bin/npm install

  ../../bin/npm run generate

  cd ../../services/frontend

  echo "Done."
}

main() {
  while [[ "$#" -gt 0 ]]; do
    case $1 in
      -bp|--build-proto)
        build_proto && exit 0;
        shift
      ;;
      *)
        echo "Unknown parameter passed: ${1}";
        exit 1
      ;;
    esac;
    shift;
  done
}

main "${@}"
