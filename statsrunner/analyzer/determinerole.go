package analyzer

import "strings"

func determineRole(lane, role string) string {
	switch strings.ToUpper(lane) {
	case "TOP":
		return "TOP"
	case "MID":
		fallthrough
	case "MIDDLE":
		return "MIDDLE"
	case "JUNGLE":
		return "JUNGLE"
	case "BOT":
		fallthrough
	case "BOTTOM":
		switch strings.ToUpper(role) {
		case "DUO_CARRY":
			return "CARRY"
		case "DUO_SUPPORT":
			return "SUPPORT"
		default:
			return "BOTTOM_UNKNOWN"
		}
	default:
		return "UNKNOWN"
	}
}
