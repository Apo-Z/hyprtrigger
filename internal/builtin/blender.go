package builtin

import "hyprtrigger/internal/events"

func registerBlender(r *events.Registry) {
	r.RegisterBuiltin(&events.Event{
		Name:     "windowtitlev2",
		Regex:    "Preferences",
		Command:  `hyprctl --batch "dispatch setfloating address:0x{WINDOW_ID}; dispatch resizewindowpixel exact 20% 50%, address:0x{WINDOW_ID}; dispatch centerwindow"`,
		UseShell: true,
	})
}
