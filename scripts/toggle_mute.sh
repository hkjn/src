#!/bin/bash

# Toggles mute/unmute of Master channel.

# First, figure out what state Master channel is in.
MASTER_OFF=$(amixer sget Master | grep -i '\[off\]');
if [ -z "$MASTER_OFF" ]; then
  echo "Master was on; will switch Master + Headphone OFF"
	amixer sset Master off
	amixer sset Headphone off
else
  echo "Master was off; will switch Master + Headphone ON"
	amixer sset Master on
	amixer sset Headphone on
fi

# TODO(hkjn): Eliminate this blunt force; amixer command above
# *claims* to have toggled sound back on, but at least on headphones
# toggling mute/unmute does not restore sound without re-init..
# IS_OFF=$(amixer sget Master | grep -i '\[off\]');
# if [ -z "$IS_OFF" ]; then
#  echo "switched back on; running alsactl init"
#  alsactl init
#fi
