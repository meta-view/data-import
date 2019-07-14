#!/bin/bash

TABLER_VERSION=0.0.33

mkdir -p assets/tabler tmp
rm -Rf tmp/tabler

wget -c -N http://maven.javastream.de/org/webjars/tabler/${TABLER_VERSION}/tabler-${TABLER_VERSION}.jar -O tmp/tabler-${TABLER_VERSION}.jar
unzip tmp/tabler-${TABLER_VERSION}.jar -d tmp/tabler
cp -R tmp/tabler/META-INF/resources/webjars/tabler/${TABLER_VERSION}/assets/ assets/tabler/
