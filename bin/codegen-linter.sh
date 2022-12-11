#!/usr/bin/env bash

set -e

function verify_hashes() {
  branch=$(tar -cf - "${path}" | md5sum)

  git checkout origin/"${GITHUB_BASE_REF}" "${path}"
  
  main=$(tar -cf - "${path}" | md5sum)

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

  modified_paths=$(git diff --name-only origin/"${GITHUB_BASE_REF}"...origin/"${GITHUB_HEAD_REF}" .)
  echo "$modified_paths"
    
  for file in "${modified_paths[@]}"; do
    
    if [[ "${file}" == *"services/migration/src/perspex"* ]]; then
      path="schemas/perspex/pkg/models"
      tool="sqlboiler"
      verify_hashes
    elif [[ "${file}" == *"schemas/graphql"* ]]; then
      path="services/backend/pkg/graphql"
      tool="gqlgen"
      verify_hashes
    elif [[ "${file}" == *" schemas/proto/**/*.proto"* ]]; then
      path="schemas/proto/goproto"
      tool="buf"
      verify_hashes
      buf_lint
    else
      echo "Generated code is up to date"
    fi
  done

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
