#!/usr/bin/env bash

set -e

path="/mnt/db"

main() {
    migration -p ${path} -m up
}

main "$@"