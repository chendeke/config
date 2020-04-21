#!/usr/bin/env bash
mkdir 1
cp config-real.yaml 1/config.yaml
ln -s 1 ..data
ln -s ..data/config.yaml config.yaml