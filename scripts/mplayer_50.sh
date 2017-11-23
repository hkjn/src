#!/bin/bash

# Starts mplayer with fixed resolution.
mplayer -vf dsize=1024:-2 $1
