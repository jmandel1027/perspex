#!/bin/bash

set -e


build_boil() {
  cd ../../schemas/perspex

  ../../bin/go work sync
 
  ../../bin/go mod download

  ../../.hermit/go/bin/sqlboiler psql
  
  cd ../../services/backend

  echo "Done."
}

build_gql() {
  cd ../../schemas/graphql

  #../../bin/go work sync
  
  #../../bin/go mod download

  #rm -rf pkg

  time ../../.hermit/go/bin/gqlgen generate

  cd ../../services/backend

  echo "Done."
}

build_proto() {
  cd ../../schemas/proto

  ../../bin/go work sync

  ../../bin/buf generate

  cd ../../services/backend

  echo "Done."
}

build_linux() {
  output="$outputPath/$app"
  src="$srcPath/$app/$pkgFile"

  echo "Building: ${app}"

  GOOS=linux GOARCH=amd64 ../../bin/go build -ldflags="-w -s" -o "${output}" "${src}"

  echo "Built: ${app} size:"

  find . -name "${output}" | awk '{print $5}'

  echo "Done building: ${app}"

  exit 0
}

run_tilt() {
  # Download go dependencies
  ../../bin/go mod download

  build_gql

  echo "Done, generating db models"
  
  build_boil
  
  echo "Done, generating proto"

  build_proto

  exit 0
}

main() {
  while [[ "$#" -gt 0 ]]; do
    case $1 in
      -sp|--source-path)
        srcPath="${2}"
        shift
        ;;
      -pf|--package-file)
        pkgFile="${2}"
        shift
        ;;
      -op|--output-path)
        outputPath="${2}"
        shift
        ;;
      -an|--app-name)
        app="${2}"
        shift
        ;;
      -bl|--build-linux)
        build_linux
        shift
      ;;
      -bb|--build-boil)
        build_boil && exit 0;
        shift
      ;;
      -bg|--build-gql)
        build_gql && exit 0;
        shift
      ;;
      -bp|--build-proto)
        build_proto && exit 0;
        shift
      ;;
      -rt|--run-tilt)
        run_tilt
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
