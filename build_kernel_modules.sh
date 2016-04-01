#!/bin/bash
function build {
    KERNEL=$1
    UNAME_R=$2
    docker build --build-arg KERNEL_VERSION=$KERNEL -t clusterhq/build-zfs-boot2docker:${UNAME_R} -f Dockerfile.boot2docker .
    UNAME_R=$UNAME_R ./zfs-builder sh -c "docker run -e UNAME_R=$UNAME_R -v ${PWD}/rootfs:/rootfs clusterhq/build-zfs-boot2docker:${UNAME_R} /build_zfs.sh && cp rootfs/zfs-${UNAME_R}.tar.gz ."
}

# travis precise, probably an ubuntu kernel
build 3.13.0 3.13.0-63-generic

# boot2docker
# look up docker version -> kernel mapping here:
# https://github.com/boot2docker/boot2docker/releases

build 4.0.9 4.0.9-boot2docker
build 4.0.10 4.0.10-boot2docker

build 4.1.12 4.1.12-boot2docker
build 4.1.13 4.1.13-boot2docker

build 4.1.17  4.1.17-boot2docker
build 4.1.18  4.1.18-boot2docker
build 4.1.19  4.1.19-boot2docker
