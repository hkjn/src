#!/bin/bash

# Deploys {inc,dec}_brightness.go.
go build inc_intel_backlight.go
go build dec_intel_backlight.go
sudo chown -v root:root {inc,dec}_intel_backlight
sudo chmod -v 4755 {inc,dec}_intel_backlight
sudo mv -v {inc,dec}_intel_backlight /usr/bin/

# Deploys power.go.
go build check_battery.go
sudo mv -v check_battery /usr/bin/

# Deploys gobatti.go.
go build gobatti.go
sudo mv -v gobatti /usr/bin/
