#!/bin/bash

mkdir mycontainer
cd mycontainer

mkdir rootfs

docker export $(docker create busybox) | tar -C rootfs -xvf -
