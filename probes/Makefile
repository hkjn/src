include ../make/Makefile
NAME=probes
VERSION=$(shell cat VERSION)
PUSH_ENABLED=0

pre-build:
	@echo "Copying in prober dependency to avoid trying to clone entire hkjn.me/src."
	cp -vr ../prober .
