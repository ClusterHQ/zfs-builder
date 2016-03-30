#!/bin/bash
#
# ZFS builder for boot2docker
#
# Needs kernel config
# SPL=y
# ZFS=y

set -xe

git clone https://github.com/zfsonlinux/spl.git /zfs/spl
cd /zfs/spl
git checkout spl-0.6.5.4

git clone https://github.com/zfsonlinux/zfs.git /zfs/zfs
cd /zfs/zfs
git checkout zfs-0.6.5.4

cd /linux-kernel
make prepare

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
./copy-builtin /linux-kernel

# Configure and cross-compile SPL usermode utils
./configure \
    --prefix=/ \
    --libdir=/lib \
    --includedir=/usr/include \
    --datarootdir=/usr/share \
    --enable-linux-builtin=yes \
    --with-linux=/linux-kernel \
    --with-linux-obj=/linux-kernel \
    --with-config=user \
    --build=x86_64-pc-linux-gnu \
    --host=x86_64-pc-linux-gnu
make
make install DESTDIR=/rootfs

# Configure and compile ZFS kernel module
cd /zfs/zfs
./autogen.sh
./configure \
    --prefix=/ \
    --libdir=/lib \
    --includedir=/usr/include \
    --datarootdir=/usr/share \
    --enable-linux-builtin=yes \
    --with-linux=/linux-kernel \
    --with-linux-obj=/linux-kernel \
    --with-spl=/zfs/spl \
    --with-spl-obj=/zfs/spl \
    --with-config=kernel
./copy-builtin /linux-kernel

# Configure and cross-compile ZFS usermode utils
./configure \
    --prefix=/ \
    --libdir=/lib \
    --includedir=/usr/include \
    --datarootdir=/usr/share \
    --with-linux=/linux-kernel \
    --with-linux-obj=/linux-kernel \
    --with-spl=/zfs/spl \
    --with-spl-obj=/zfs/spl \
    --with-config=user \
    --build=x86_64-pc-linux-gnu \
    --host=x86_64-pc-linux-gnu
make
make install DESTDIR=/rootfs
