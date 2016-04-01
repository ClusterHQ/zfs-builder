#!/bin/bash
docker build -t clusterhq/build-zfs-boot2docker:docker-1.8.1-linux-4.0.9 -f Dockerfile.boot2docker .
# TODO make the container output the binaries somewhere...
./zfs-builder docker run -v ${PWD}/rootfs:/rootfs clusterhq/build-zfs-boot2docker /build_zfs.sh
