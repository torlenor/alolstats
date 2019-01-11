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

	// Specifies Summoner names for which matches shall be fetched
	FetchMatchesForSummoners []string
	// How many of the last matches shall be checked/pulled per account. 0 means all of them
	FetchMatchesForSummonersNumber uint64

	// Specified for which leagues matches shall be fetched. Currently implemented by Riot are "masterleagues", "grandmasterleagues", "challengerleagues"
	FetchMatchesForLeagues []string
	// Specified for queues matches shall be fetched. Allowed are "RANKED_SOLO_5x5", "RANKED_FLEX_SR", "RANKED_FLEX_TT"
	FetchMatchesForLeagueQueues []string
	// How many of the last matches shall be checked/pulled per account. 0 means all of them
	FetchMatchesForLeaguesNumber uint64
}

// StatsRunner holds the settings for the StatsRunner
type StatsRunner struct {
	RunRScripts      bool   // Specifies if R scripts shall be used (needs a running R installation)
	RScriptPath      string // Path to the R scripts (distributed with alolstats)
	RPlotsOutputPath string // Path where the generated plots shall be stored

	RScriptsUpdateInterval uint32 // Update Interval for running the R scripts in minutes > 0

	GameVersion string // We want to do stats calculations for the following versions, must be valid game versions, see https://ddragon.leagueoflegends.com/api/versions.json, e.g., 9.1.1, 8.24.1
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
