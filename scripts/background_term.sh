#!/bin/bash

# Creates a xfce4-terminal that's always in fullscreen in the back,
# without menus or other widgets.

cd $HOME

xfce4-terminal --hide-borders --hide-toolbar --hide-menubar --title=desktopconsole --geometry=192x58+0+0 &

sleep 1
wmctrl -r desktopconsole -b add,below,sticky
wmctrl -r desktopconsole -b add,skip_taskbar,skip_pager

