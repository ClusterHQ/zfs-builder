#!/bin/bash
docker build -t clusterhq/build-zfs-boot2docker -f Dockerfile.boot2docker .
# TODO make the container output the binaries somewhere...
./zfs-builder docker run -v ${PWD}:/pwd clusterhq/build-zfs-boot2docker
