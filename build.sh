#!/bin/bash

go test ./...
go build --tags "icu json1 fts5 secure_delete" -o bin/meta-view-service

