#!/bin/bash

TABLER_VERSION=0.0.33
DROPZONE_VERSION=5.5.0
DROPZONE_BUILD_JOBID=79974233 #5.5.0

rm -Rf tmp/tabler assets/tabler
mkdir -p assets/tabler/js assets/js tmp/tabler

wget -c http://maven.javastream.de/org/webjars/tabler/${TABLER_VERSION}/tabler-${TABLER_VERSION}.jar -O tmp/tabler-${TABLER_VERSION}.jar
unzip tmp/tabler-${TABLER_VERSION}.jar -d tmp/tabler
mv tmp/tabler/META-INF/resources/webjars/tabler/${TABLER_VERSION}/assets/* assets/tabler/
mv assets/tabler/js/core.js assets/js/
mv -f assets/tabler/js/vendors assets/js/

wget -c https://gitlab.com/meno/dropzone/-/jobs/${DROPZONE_BUILD_JOBID}/artifacts/raw/dist/dropzone.js -O tmp/dropzone-${DROPZONE_VERSION}.js
mv tmp/dropzone-${DROPZONE_VERSION}.js assets/js/dropzone.js

rm -Rf tmp/*

go get -u github.com/jteeuwen/go-bindata/...
go-bindata -o assets/assets.go assets/...
sed -i'.bak' -e 's/package main/package assets/g' assets/assets.go && rm assets/assets.go.bak