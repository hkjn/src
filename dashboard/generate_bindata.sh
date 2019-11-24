#!/usr/bin/env bash


go-bindata -pkg gen -o gen/bindata.go probes.yaml tmpl/
