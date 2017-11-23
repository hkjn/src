#!/bin/bash

# Raises volume of Master + Headphone channel slightly.
amixer -c 0 set Master 2.5dB+
amixer -c 0 set Headphone 2.5dB+
