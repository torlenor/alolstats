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

// MemoryBackend holds the settings for the memory backend
type MemoryBackend struct {
}

// MongoBackend holds the settings for the memory backend
type MongoBackend struct {
	// URL to connect to
	URL string
	// Database to use
	Database string
}

// StorageBackend holds the settings for the used backend component
type StorageBackend struct {
	// Name of the storage backend to use (e.g., memory, mongo)
	Backend string

	MemoryBackend MemoryBackend
	MongoBackend  MongoBackend
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

// FetchRunner holds the settings for the FetchRunner
type FetchRunner struct {
	// Specified the update interval for fetching Summoner Matches in minutes > 0
	UpdateIntervalSummonerMatches uint32
	// Specifies AccountIDs for Summoners where Matches shall be fetched
	MatchesForSummonerAccountIDs []uint64
}

// StatsRunner holds the settings for the StatsRunner
type StatsRunner struct {
}

// Config holds the complete ALolStats config
type Config struct {
	API            API
	RiotClient     RiotClient
	LoLStorage     LoLStorage
	StorageBackend StorageBackend
	FetchRunner    FetchRunner
	StatsRunner    StatsRunner
}
