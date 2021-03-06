#!/bin/bash

if [[ ! -d "assets" ]]; then
    echo ""
    echo "assets folder does not exist. run update.sh"
    echo ""
    exit 1
fi

go test ./...
go build \
    -ldflags "-X main.VersionString=$(git describe --always --dirty --tags)" \
    --tags "icu json1 fts5 secure_delete" -o bin/meta-view-service

