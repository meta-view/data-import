#!/bin/bash

TABLER_VERSION=0.0.33

mkdir -p assets/tabler assets/js tmp
rm -Rf tmp/tabler assets/tabler

wget -c -N http://maven.javastream.de/org/webjars/tabler/${TABLER_VERSION}/tabler-${TABLER_VERSION}.jar -O tmp/tabler-${TABLER_VERSION}.jar
unzip tmp/tabler-${TABLER_VERSION}.jar -d tmp/tabler
cp -R tmp/tabler/META-INF/resources/webjars/tabler/${TABLER_VERSION}/assets/ assets/tabler/
mv assets/tabler/js/core.js assets/js/
mv -f assets/tabler/js/vendors assets/js/