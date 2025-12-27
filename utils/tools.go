package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func SafeString(data bson.M, key string) string {

	if data == nil {
		return ""
	}

	val, ok := data[key]
	if !ok || val == nil {
		return ""
	}

	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}

// Helper to create Clickable Links (OSC 8)
func MakeLink(styledText, url string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, styledText)
}
