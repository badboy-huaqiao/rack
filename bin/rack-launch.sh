#!/bin/bash

DIR=$PWD
CMD=../cmd/rack
function cleanup {
	pkill rack
}

cd $CMD
exec -a rack ./rack &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done