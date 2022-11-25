#!/usr/bin/env bash

# Getting sqlboiler/gqlgen dependencies
go get \
    github.com/volatiletech/sqlboiler \
    github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql \
    github.com/cortesi/modd/cmd/modd \
    github.com/vektah/dataloaden \
    github.com/99designs/gqlgen

# Download go dependencies
go mod download

echo "Checking if $(pwd)/schema exists and creating if not"
if [ ! -d schema ]; then
    echo "Creating $(pwd)/schema"
    mkdir -p schema
else
    echo "$(pwd)/schema exists, emptying it"
    rm -rf schema/*
fi

echo "Grabbing all schema files served using the python http.server (localhost:8000)"
cd schema && wget -q -r -np -nH --cut-dirs=3 -R index.html http://localhost:8000

echo "The following schema files are now present:"
ls
cd ..

echo "Running '/bin/prepare -bg'"
bin/prepare -bg
