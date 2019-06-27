#!/bin/bash

TABLER_VERSION=0.0.32

mkdir -p assets/tabler tmp
rm -Rf tmp/tabler

wget -c -N https://repo1.maven.org/maven2/org/webjars/tabler/${TABLER_VERSION}/tabler-${TABLER_VERSION}.jar -O tmp/tabler-${TABLER_VERSION}.jar
unzip tmp/tabler-${TABLER_VERSION}.jar -d tmp/tabler
cp -R tmp/tabler/META-INF/resources/webjars/tabler/${TABLER_VERSION}/assets/ assets/tabler/
