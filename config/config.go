package config

// API holds the API settings for the AbylALolStatseBotter configuration
type API struct {
	// IP Address the REST API listens on
	// If empty or non-existing listen on all interfaces
	IP string
	// Port the REST API listens on
	Port string
}

// RiotClient holds the settings specific for the Riot API
type RiotClient struct {
	// Riot developer API key used for API access
	Key string
	// Riot API version (v3, v4, ...)
	APIVersion string
	// Game region to use ("euw1", "eun1", ...)
	Region string
}

// LoLStorage holds the settings specific for the storage component
type LoLStorage struct {
	// Name of the storage backend to use (e.g., sqlite)
	Backend string
	// Specifies if Riot provided Match Files should be read
	UseMatchFiles bool
	// Specifies the directory holding the match files
	MatchFileDir string
	// Specified the maximum age for champion data in minutes until it's invalidated. 0 means it is always fetched newly.
	MaxAgeChampion uint32
	// Specified the maximum age for free champion rotation data in minutes until it's invalidated. 0 means it is always fetched newly.
	MaxAgeChampionRotation uint32
	// Specified the maximum age for summoner data in minutes until it's invalidated. 0 means it is always fetched newly.
	MaxAgeSummoner uint32
}

// Config holds the complete ALolStats config
type Config struct {
	API        API
	RiotClient RiotClient
	LoLStorage LoLStorage
}
