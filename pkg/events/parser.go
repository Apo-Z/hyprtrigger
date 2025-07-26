package events

import "strings"

func ParseEventData(eventName, rawData string) *EventData {
	switch eventName {
	case "windowtitlev2":
		parts := strings.SplitN(rawData, ",", 2)
		if len(parts) >= 2 {
			return &EventData{
				WindowID: parts[0],
				Content:  parts[1],
			}
		}
	case "openwindow":
		parts := strings.SplitN(rawData, ",", 4)
		if len(parts) >= 4 {
			return &EventData{
				WindowID: parts[0],
				Content:  parts[3],
			}
		}
	case "activewindow":
		// Format attendu: "windowID,windowTitle"
		if commaIndex := strings.Index(rawData, ","); commaIndex != -1 {
			return &EventData{
				WindowID: strings.TrimSpace(rawData[:commaIndex]),
				Content:  strings.TrimSpace(rawData[commaIndex+1:]),
			}
		} else {
			// Pas de virgule trouvée, tout est considéré comme WindowID
			return &EventData{
				WindowID: strings.TrimSpace(rawData),
				Content:  "",
			}
		}
	}

	return &EventData{
		WindowID: "",
		Content:  rawData,
	}
}
