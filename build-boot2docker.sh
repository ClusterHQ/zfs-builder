#!/bin/bash
function build {
    KERNEL=$1
    DOCKER=$2
    docker build --build-arg KERNEL_VERSION=$KERNEL -t clusterhq/build-zfs-boot2docker:docker-${DOCKER}-linux-${LINUX} -f Dockerfile.boot2docker .
    # TODO make the container output the binaries somewhere...
    ./zfs-builder docker run -e KERNEL_VERSION=$KERNEL_VERSION -v ${PWD}/rootfs:/rootfs clusterhq/build-zfs-boot2docker /build_zfs.sh ${KERNEL}
    cp rootfs/zfs-$KERNEL_VERSION.tar.gz .
}

# look up docker version -> kernel mapping here:
# https://github.com/boot2docker/boot2docker/releases

build 4.0.9  1.8.2
build 4.0.10 1.8.3

build 4.1.12 1.9.0
build 4.1.13 1.9.1

build 4.1.17 1.10.1
build 4.1.18 1.10.2
build 4.1.19 1.10.3
