package utils

import (
	"fmt"
	"log"
	"portfolioTUI/config"
	"portfolioTUI/database"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// GetIcon returns a specific Nerd Font symbol based on the category string
func GetIcon(category string) string {
	// Normalize string to lowercase for easier matching
	c := strings.ToLower(category)

	if strings.Contains(c, "web") {
		return "" // Globe (Web Development)
	}
	if strings.Contains(c, "app") || strings.Contains(c, "mobile") {
		return "" // Mobile Phone (App Development)
	}
	if strings.Contains(c, "cms") || strings.Contains(c, "wordpress") {
		return "" // Box / Package (CMS Solutions)
	}
	if strings.Contains(c, "video") || strings.Contains(c, "edit") {
		return "" // Film Strip (Video Editing)
	}
	if strings.Contains(c, "cloud") || strings.Contains(c, "devops") {
		return "" // Cloud Icon (Cloud Solutions)
	}
	if strings.Contains(c, "backend") || strings.Contains(c, "api") {
		return "" // Database (Backend)
	}
	if strings.Contains(c, "design") || strings.Contains(c, "ui") {
		return "" // Paint Brush (Design)
	}

	return "🛠️" // Default Tool
}

func SafeID(data map[string]interface{}, key string) string {
	val, ok := data[key]
	if !ok || val == nil {
		return ""
	}

	// 1. BEST CASE: It is a primitive.ObjectID
	if oid, ok := val.(primitive.ObjectID); ok {
		return oid.Hex() // Returns just "6829..."
	}

	// 2. FALLBACK: It's already a string or something else
	// Convert to string first
	str := fmt.Sprintf("%v", val)

	// 3. CLEANUP: Strip "ObjectID(" and ")" if they exist
	// This removes 'ObjectID("' from the start
	if strings.HasPrefix(str, "ObjectID(\"") {
		str = strings.TrimPrefix(str, "ObjectID(\"")
		str = strings.TrimSuffix(str, "\")")
	} else if strings.HasPrefix(str, "ObjectID(") {
		// Handle case without quotes just in case: ObjectID(123)
		str = strings.TrimPrefix(str, "ObjectID(")
		str = strings.TrimSuffix(str, ")")
	}

	return str
}

// SafeDate converts timestamps (ms), ISO strings, or string-numbers to a readable format
func SafeDate(data map[string]interface{}, key string) string {
	val, ok := data[key]
	if !ok || val == nil {
		return "No Date"
	}

	var t time.Time

	// DEBUG PRINT: This will show up in your terminal
	// fmt.Printf("DEBUG DATE [%s]: Type=%T Value=%v\n", key, val, val)

	switch v := val.(type) {

	// Case 1: MongoDB DateTime
	case primitive.DateTime:
		t = v.Time()

	// Case 2: Standard Float64 (JSON numbers often come as this)
	case float64:
		t = time.UnixMilli(int64(v))

	// Case 3: Standard Int64
	case int64:
		t = time.UnixMilli(v)

	// Case 4: String containing a number ("1748274988656")
	case string:
		// Try parsing as integer number
		if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
			t = time.UnixMilli(ms)
		} else {
			// Try parsing as ISO Date ("2025-05-26T...")
			parsed, err := time.Parse(time.RFC3339, v)
			if err == nil {
				t = parsed
			} else {
				return v // Failed to parse, return raw string
			}
		}

	default:
		// If we get here, the type is unknown.
		// Converting whatever it is (likely int32 or uint) to int64 might fix it.
		// Let's try to force it to a string and parse that.
		strVal := fmt.Sprintf("%v", v)
		if ms, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			t = time.UnixMilli(ms)
		} else {
			return strVal
		}
	}

	return t.Format("Jan 02, 2006")
}

func FetchData() tea.Cmd {
	return func() tea.Msg {
		var alldata config.AllMessages
		// Use your new generic function
		for _, collections := range config.Collection { // Ensure config.Collections matches your config file
			data, err := database.GetALLFromCollection(collections, collections)
			if err != nil {
				log.Println("Error fetching", collections, err)
				continue
			}
			log.Println("Fetched", len(data), "items from", collections)
			alldata = append(alldata, config.DataMsg{Type: collections, Data: data})
		}
		return alldata
	}
}

// submitContactCmd unpacks the form struct and calls the DB function
func submitContactCmd(form database.ContactSchema) tea.Cmd {
	return func() tea.Msg {
		// 1. Simulate Delay
		time.Sleep(1 * time.Second)

		// 2. Handle the ServiceID pointer safely
		// The struct has *string, but InsertContact expects a string ID
		svcID := ""
		if form.ServiceId != nil {
			svcID = *form.ServiceId
		}

		// 3. Call the DB function with individual fields
		err := database.InsertContact(
			form.FirstName,
			form.LastName,
			form.Email,
			form.Type,
			form.Description,
			svcID,
		)

		if err != nil {
			return config.FormSubmittedMsg{Success: false}
		}
		return config.FormSubmittedMsg{Success: true}
	}
}

func GenerateImagesCmds(dataType string, data []bson.M) []tea.Cmd {
	var cmds []tea.Cmd
	defaultImage := config.DEFAULTIMAGEURL

	for i, item := range data {
		var url string
		var key string
		var targetWidth int
		var targetHeight int

		switch dataType {
		case "projects":
			targetWidth = 32
			targetHeight = 15
		case "positions":
			targetWidth = 30
			targetHeight = 15
		case "blogs":
			targetWidth = 40
			targetHeight = 12
		default:
			targetWidth = 30
			targetHeight = 15

		}

		// 1. Determine the correct key for this collection
		switch dataType {
		case "projects":
			key = "imageUrl"
		case "positions":
			key = "logoUrl" // Ensure this matches your DB
		case "blogs":
			key = "featuredImage" // Ensure this matches your DB
		default:
			continue
		}

		// 2. Safe Get: Get the string, or "" if missing
		if val, ok := item[key].(string); ok {
			url = val
		}

		// 3. Fallback Logic (The Fix)
		// If URL is empty OR too short, use the default
		if len(url) < 5 {
			url = defaultImage
		}

		// 4. Generate the Command
		// Now we always have a valid URL (either original or default)
		cmds = append(cmds, GenerateAsciiImage(url, dataType, i, targetWidth, targetHeight))
	}
	return cmds
}
