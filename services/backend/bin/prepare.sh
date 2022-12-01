#!/bin/bash

set -e


build_boil() {
  cd ../../schemas/perspex
  
  go mod download

  go get \
    github.com/volatiletech/sqlboiler \
    github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql

  sqlboiler psql
  
  cd ../../services/backend

  printf "\nDone.\n\n"
}

build_gql() {
  cd ../../schemas/graphql

  go mod download

  go get github.com/99designs/gqlgen

  rm -rf pkg/resolvers/generated
  rm -rf pkg/graphql/*

  time go run -v github.com/99designs/gqlgen generate
  time go generate ./...

  cd ../../services/backend

  printf "\nDone.\n\n"
}

build_linux() {
  output="$outputPath/$app"
  src="$srcPath/$app/$pkgFile"

  echo "Building: ${app}"
  GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o "${output}" "${src}"
  echo "Built: ${app} size:"
  ls -lah "${output}" | awk '{print $5}'
  echo "Done building: ${app}"
  exit 0
}

build_mac() {
  output="$outputPath/$app"
  src="$srcPath/$app/$pkgFile"

  printf "\nBuilding: $app\n"
  go build -o $output $src
  printf "\nBuilt: $app size:"
  ls -lah $output | awk '{print $5}'
  printf "\nDone building: $app\n\n"
  exit 0
}


run_linux() {
  buildPath="bin"
  app="server"
  program="$buildPath/$app"
  printf "\nStart app: $app\n"
  printenv

  # Set all ENV vars for the program to run
  # export $(grep -v '^#' ./.env | xargs)
  $program

  # This should unset all the ENV vars, just in case.
  # unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)
  printf "\nStopped app: $app\n\n"
  exit 0
}

run_mac() {
    app="server"
    src="$srcPath/$app/$pkgFile"
    printf "\nStart running: $app\n"

    time modd
    # This should unset all the ENV vars, just in case.
    # unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/''' | xargs)
    printf "\nStopped running: $app\n\n"
    exit 0
}

run_tilt() {

  # Download go dependencies
  go mod download

  build_gql

  echo "Done, generating db models"

  build_boil
  
  exit 0
}

main() {
  while [[ "$#" -gt 0 ]]; do
    case $1 in
      -sp|--source-path)
        srcPath="$2"
        shift
        ;;
      -pf|--package-file)
        pkgFile="$2"
        shift
        ;;
      -op|--output-path)
        outputPath="$2"
        shift
        ;;
      -an|--app-name)
        app="$2"
        shift
        ;;
      -bl|--build-linux)
        build_linux
        shift
      ;;
      -bb|--build-boil)
        build_boil
        shift
      ;;
      -bg|--build-gql)
        build_gql
        shift
      ;;
      -bm|--build-mac)
        build_mac
        shift
      ;;
      -rl|--run-linux)
        run_linux
        shift
      ;;
      -rt|--run-tilt)
        run_tilt
        shift
      ;;
      -rm|--run-mac)
        run_mac
        shift
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