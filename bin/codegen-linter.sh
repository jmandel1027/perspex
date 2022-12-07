#!/usr/bin/env bash

set -e

function verify_hashes() {
  echo "about to verify hashes"
  branch=$(tar -cf - "${genpath}" | md5sum)
  echo "branch: ${branch}"

  changes=$(git pull origin main --quiet || echo $?)
  echo "changes: ${changes}"
  if [[ "${changes}" == "Already up to date." ]]; then
    echo "Error: Generated ${tool} code is out of phase, please commit generated code."
    exit 1;
  fi

  echo "about to checkout ${genpath} main"
  git checkout --theirs "${genpath}" --quiet main
  
  main=$(tar -cf - "${genpath}" | md5sum)
  echo "main: ${main}"

  if [[ "${branch}" == "${main}" ]]; then
    echo "Error: Generated ${tool} code is out of phase, please commit generated code."
    exit 1;
  fi

  echo "about to run buf lint"
  if [[ "${tool}" == "buf" ]]; then
    buf_lint
  fi

  echo "done"
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
