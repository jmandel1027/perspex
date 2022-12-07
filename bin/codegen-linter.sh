#!/usr/bin/env bash

set -e

function verify_hashes() {
  branch=$(tar -cf - "${genpath}" | md5sum)
  
  git checkout main "${genpath}"
  
  main=$(tar -cf - "${genpath}" | md5sum)

  if [[ "${branch}" == "${main}" ]]; then
    echo "Error: Generated ${tool} code is out of phase, please commit generated code."
    exit 1;
  fi

  if [[ "${tool}" == "buf" ]]; then
    buf_lint
  fi
}

function buf_lint() {
  cd schemas/proto
  ../../bin/buf lint
  cd ../..
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
       -g|--genpath)
        genpath="$2"
        shift
        ;;
      -t|--tool)
        tool="$2"
        shift
        ;;
      -l | --lint )
        verify_hashes
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
