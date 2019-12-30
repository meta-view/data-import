#!/bin/bash

 docker run -v $PWD/data:/data -p 9000:9000 phaus/meta-view-service:latest
