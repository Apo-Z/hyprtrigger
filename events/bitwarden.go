package events

import "hyprtrigger/pkg/events"

func RegisterBitwardenEvents() {
	events.RegisterEvent(&events.Event{
		Name:     "windowtitlev2",
		Regex:    "Bitwarden",
		Command:  `hyprctl --batch "dispatch setfloating address:0x{WINDOW_ID}; dispatch resizewindowpixel exact 20% 50%, address:0x{WINDOW_ID}; dispatch centerwindow"`,
		UseShell: true,
	})
}

func init() {
	RegisterBitwardenEvents()
}
