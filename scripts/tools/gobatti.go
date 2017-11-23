// Simple tool to show battery status from SysFS as GTK icon.
//
// Icons are taken from the current GTK theme, typically stored at
// /usr/share/icons/gnome/32x32/status or similar.
package main

import (
	"log"
	"os"
	"time"

	"github.com/mattn/go-gtk/gtk"
	"hkjn.me/power"
)

var levels = []string{
	"empty",
	"caution",
	"low",
	"good",
	"full",
}

var states = []string{
	"charged",
	"charging",
}

type Icon struct {
	Battery    power.Battery
	StatusIcon *gtk.StatusIcon
}

func getIcons() []Icon {
	bat, err := power.Get()
	if err != nil {
		log.Fatalf("failed to get battery info: %v\n", err)
	}
	result := make([]Icon, len(bat))
	for i, b := range bat {
		icon := Icon{
			Battery: b,
		}
		icon.create()
		result[i] = icon
	}
	return result
}

func (i *Icon) create() {
	in := i.Battery.Desc()
	i.StatusIcon = gtk.NewStatusIconFromIconName(in)
	i.StatusIcon.SetTooltipText(i.Battery.String())
	log.Printf("Created status icon %v with icon name %s\n", i, in)
}

// update updates the icon with new battery info.
func (i *Icon) update(battery power.Battery) {
	oldName := i.Battery.Desc()
	i.Battery = battery
	newName := i.Battery.Desc()
	if newName != oldName {
		log.Printf("Changing icon to %q from %q..\n", newName, oldName)
		// TODO: this should check if the icon name to set actually exists, somehow.
		i.StatusIcon.SetFromIconName(newName)
	}
	i.StatusIcon.SetTooltipText(battery.String())
}

// poll reads info from battery i and sleeps for specified duration.
func poll(d time.Duration, i int, icon Icon) {
	for {
		b, err := power.GetNumber(i)
		if err != nil {
			if err == power.ErrNoFile || err == power.ErrNoDevice {
				// No SysFS file; set to unknown and assume that battery
				// returns eventually. This happens i.e. if battery is
				// physically disconnected, or just after resuming from being
				// suspended.
				// TODO: possibly give up and drop icon after N tries?
				log.Printf("no sysfs file for battery %d, setting state to Unknown\n", i)
				b.State = power.Unknown
			} else {
				log.Fatalf("failed to get battery info for battery %d: %v\n", i, err)
			}
		}
		log.Printf("[Battery %d]: %+v\n", i, b)
		icon.update(b)
		time.Sleep(d)
	}
}

func main() {
	gtk.Init(&os.Args)
	// TODO: we currently assume that number of batteries at startup
	// never can grow (but can shrink).
	icons := getIcons()
	d, err := time.ParseDuration("20s")
	if err != nil {
		log.Fatalf("bad duration: %v\n", err)
	}
	for i := 0; i < len(icons); i++ {
		go poll(d, i, icons[i])
	}

	log.Printf("Calling gtk.Main()..")
	gtk.Main()
}
