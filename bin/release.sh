#!/usr/bin/env bash

set -e

function disable_services() {
  services=$(find services -maxdepth 1 -type d)
  
  names=$()
  for service in "${services[@]}"; do
    names=$(echo "${service}" | cut -d'/' -f2-)
  done

  read -r -a service_array <<< ${names}

  for s in "${service_array[@]}"; do
    if [[ "${s}" == "services" ]]; then
      continue
    fi
    
    tilt_var=$(echo "TILT_${s}_ENABLED" | tr "[:lower:]" "[:upper:]")
    if ! [[ "${local}" == "true" ]]; then
      echo "${tilt_var}=false"
    else
      echo "${tilt_var}=false" >> "${GITHUB_ENV}" 
    fi
  done

  exit 0;
}

function changed_services() {
  if ! [[ "${local}" == "true" ]]; then
    current_branch=$(git rev-parse --abbrev-ref HEAD)
    modified_services=$(git diff --name-only origin/main...origin/"${current_branch}" services \
      | grep -E "^(services)" \
      | awk -F "/" '{print $1"/"$2}' \
      | sort -u \
      | tr '\n' ' ' \
      | tr -s ' ')
  else
    modified_services=$(git diff --name-only origin/"${GITHUB_BASE_REF}" -- origin/"${GITHUB_HEAD_REF}" services \
      | grep -E "^(services)" \
      | awk -F "/" '{print $1"/"$2}' \
      | sort -u \
      | tr '\n' ' ' \
      | tr -s ' ')
  fi
 
  for service in "${modified_services[@]}"; do
    name=$(echo "${service}" | cut -d'/' -f2- | sed -e 's/\ *$//g')     
    service_name=$(echo "${name}" | tr '[:lower:]' '[:upper:]')
    if ! [[ "${local}" == "true" ]]; then
      echo "TILT_${service_name}_ENABLED=true"
      echo "TILT_${service_name}_IMAGE_TARGET=deployable" 
      bin/yq eval -i ".${service}.enabled = ${TILT_BACKEND_ENABLED}" infrastructure/tilt/values-dev.yaml
    else
      echo "TILT_${service_name}_ENABLED=true" >> "${GITHUB_ENV}"
      echo "TILT_${service_name}_IMAGE_TARGET=deployable"  >> "${GITHUB_ENV}"
      bin/yq eval -i ".${service}.enabled = ${TILT_BACKEND_ENABLED}" infrastructure/tilt/values.yaml
    fi
  done

  exit 0;
}

function display_help() {
    echo "This script manages our release processes for images and helm charts."
    echo ""
    echo "Parameters:"
    echo "    -l                      Run in local mode."
    echo "    -ds                     Disable All Services."
    echo "    -ds                     Filter for changed Services."
    echo "    -h                      Display this help message."
    echo "Usage:"
    echo "   bin/release.sh -l -ds"
    echo "   bin/release.sh -ds"
    echo "   bin/release.sh -l -s"
    echo "   bin/release.sh -s"
    echo "   bin/release.sh -h"
}

main() {
  while [[ "$#" -gt 0 ]]; do 
    case $1 in
      -l | --local )
        local=true
        ;;
      -ds | --disable-services )
        disable_services
        ;;
      -s | --services )
        changed_services
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
