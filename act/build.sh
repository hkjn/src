#!/usr/bin/env bash

docker build --build-arg FROM_IMAGE=ubuntu --build-arg FROM_TAG=22.04 .
