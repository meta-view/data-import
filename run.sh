#!/bin/bash

go run \
    -ldflags "-X main.VersionString=$(git describe --always --dirty --tags)" \
    --tags="icu json1 fts5 secure_delete" main.go

