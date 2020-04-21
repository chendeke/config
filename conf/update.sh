#!/usr/bin/env bash

NAME=20
LAST=`expr $NAME - 1`
mkdir $NAME
cp config-real.yaml $NAME/config.yaml
rm -fr ..data
rm -fr $LAST
ln -s $NAME ..data