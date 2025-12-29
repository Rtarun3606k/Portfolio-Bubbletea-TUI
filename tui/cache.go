package tui

import (
	"log"
	"portfolioTUI/config"
	"portfolioTUI/database"
	"sync"
	"time"
)

// Global Cache Storage
var (
	globalCache   config.AllMessages
	lastFetched   time.Time
	cacheMutex    sync.Mutex
	cacheDuration = 5 * time.Minute
)

// GetOrFetchData returns cached data if fresh, or fetches new data if expired
func GetOrFetchData() config.AllMessages {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 1. Check if cache is valid
	if !lastFetched.IsZero() && time.Since(lastFetched) < cacheDuration {
		log.Println("⚡ Using Cached Data (Expires in", cacheDuration-time.Since(lastFetched), ")")
		return globalCache
	}

	// 2. Cache expired (or empty), fetch fresh data
	log.Println("🔄 Cache expired or empty. Fetching fresh data from MongoDB...")
	newData := fetchAllDataSync()

	// 3. Update Cache
	if len(newData) > 0 {
		globalCache = newData
		lastFetched = time.Now()
	}

	return globalCache
}

// fetchAllDataSync gets data synchronously (blocking) for the initial load
func fetchAllDataSync() config.AllMessages {
	var alldata config.AllMessages

	// Iterate over collections from config (projects, experience, etc.)
	for _, collectionName := range config.Collection {
		data, err := database.GetALLFromCollection(collectionName, collectionName)
		if err != nil {
			log.Println("Error fetching", collectionName, err)
			continue
		}
		// Append to our list
		alldata = append(alldata, config.DataMsg{Type: collectionName, Data: data})
	}

	return alldata
}
