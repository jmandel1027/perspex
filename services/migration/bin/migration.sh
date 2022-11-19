#!/usr/bin/env bash

set -e

# change this when we implement a proper secrets integration flow w/ secrets-manager
POSTGRES_PASSWORD="pass"

dsn="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

function display_help() {
    echo "This script will run a database migration service"
    echo ""
    echo "Parameters:"
    echo "    -h                      Display this help message."
    echo "    -p                      Specifies the migration paths"
    echo "    -l                      Toggles the dsn eval for local dev, upper envs are dynamic by vars"
    echo "    -m                      Input arg for the desired commands: up|down"
    echo "    -n                      Input to specify the desired steps to increment or decrement the chain"
    echo "    -d                      Drops the entire database"
    echo "Usage:"
    echo "    bin/migration.sh -p db -l -m up"
}

create_migration() {
    migrate create -ext sql -dir ${path} -format "unix" change_me
    exit 0
}

derive_version(){
    echo $dsn
    migrate -verbose -path ${path} -database ${dsn} version
    exit 0
}

migration() {
    echo "Starting DB migration: ${cmd} ${number}"
    migrate -verbose -path ${path} -database ${dsn} ${cmd} ${number}
    exit 0
}

drop() {
    migrate -verbose -path ${path} -database ${dsn} drop -f
    exit 0
}

main() {
    while [[ "$#" -gt 0 ]]; do 
        case $1 in
            -p | --path )
                path="$2"
                shift
                ;;
            -c | --create )
                create_migration
                shift
                ;;
            -l | --local )
                dsn="postgresql://postgres:pass@127.0.0.1:5433/pg?sslmode=disable"
                ;;
            -n | --num )
                number="$2"
                shift
                ;;
            -m | --command )
                cmd="$2"
                migration
                shift
                ;;
            -d | --drop )
                drop
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

main "$@"