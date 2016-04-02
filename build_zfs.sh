#!/bin/bash
#
# ZFS builder for boot2docker

set -xe

cd /rootfs
rm -rf *

git clone https://github.com/zfsonlinux/spl.git /zfs/spl
cd /zfs/spl
git checkout spl-0.6.5.4

git clone https://github.com/zfsonlinux/zfs.git /zfs/zfs
cd /zfs/zfs
git checkout zfs-0.6.5.4

# Configure and compile SPL kernel module
cd /zfs/spl
./autogen.sh
./configure \
    --prefix=/ \
    --libdir=/lib \
    --includedir=/usr/include \
    --datarootdir=/usr/share \
    --with-linux=/linux-kernel \
    --with-linux-obj=/linux-kernel \
    --with-config=kernel

make -j8
make install DESTDIR=/rootfs

# Configure and compile ZFS kernel module
cd /zfs/zfs
./autogen.sh
./configure \
    --prefix=/ \
    --libdir=/lib \
    --includedir=/usr/include \
    --datarootdir=/usr/share \
    --with-linux=/linux-kernel \
    --with-linux-obj=/linux-kernel \
    --with-spl=/zfs/spl \
    --with-spl-obj=/zfs/spl \
    --with-config=kernel

make -j8
echo "Got after make $?"
make install DESTDIR=/rootfs
echo "Got after make install $?"

cd /rootfs/lib/modules/*
tar cfv /rootfs/zfs-${UNAME_R}.tar.gz .
