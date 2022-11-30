#!/usr/bin/env bash

set -e

function verify_sqlboiler_hashes() {
  echo "Error: Generated sqlboiler code is out of phase, please commit generated code."
  exit 1;
}

function verify_gqlgen_hashes() {
  # make hash of generated backend gqlgen files
  #~  regenerate gqlgen
  #~  check if hash is different,  if different, throw err  
  echo "Error: Generated gqlgen code is out of phase, please commit generated code."
  exit 1;
}

function lint_codegen() {
  if [[ "$(git diff --quiet HEAD main -- services/migration/src/perspex || echo $?)" == 1 ]]; then
    verify_sqlboiler_hashes
  elif [[ "$(git diff --quiet HEAD main -- schemas/graphql || echo $?)" == 1 ]]; then
    verify_gqlgen_hashes
  else
    echo "Generated code is up to date"
  fi

  exit 0;
}

function display_help() {
    echo "This script will run a database migration validation service"
    echo ""
    echo "Parameters:"
    echo "    -l | --lint             Lint to verify codegen changes have been committed."
    echo "    -h                      Display this help message."
    echo "Usage:"
    echo "   bin/codegen-linter.sh -l"
    echo "   bin/codegen-linter.sh -h"
}

main() {
  while [[ "$#" -gt 0 ]]; do 
    case $1 in
      -l | --lint )
        lint_codegen
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
