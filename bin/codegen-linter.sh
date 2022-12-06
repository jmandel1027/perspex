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
}

function buf_lint() {
  cd schemas/proto
  ../../bin/buf lint
  cd ../..
}

function lint_codegen() {
  echo "about to lint codegen"
  main_branch=$(git rev-parse --short main)
  echo $(git diff --quiet ${main_branch} -- services/migration/src/perspex || echo $?)
  if [[ "$(git diff --quiet  ${main_branch} -- services/migration/src/perspex || echo $?)" == "1" ]]; then
    echo "perspex changed"
    genpath="schemas/perspex/pkg/models"
    tool="sqlboiler"
    verify_hashes
  echo $(git diff --quiet ${main_branch} -- schemas/graphql || echo $?)
  elif [[ "$(git diff --quiet ${main_branch} -- schemas/graphql || echo $?)" == "1" ]]; then
    echo "graphql changed"
    genpath="schemas/graphql/pkg"
    tool="gqlgen"
    verify_hashes
  echo $(git diff --quiet ${main_branch} -- schemas/proto/**/*.proto || echo $?)
  elif [[ "$(git diff --quiet ${main_branch} -- schemas/proto/**/*.proto || echo $?)" == "1" ]]; then
    echo "proto changed"
    genpath="schemas/proto/goproto"
    tool="buf"
    verify_hashes
    buf_lint
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
