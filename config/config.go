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

// MongoBackend holds the settings for the mongodb backend
type MongoBackend struct {
	// URL to connect to
	URL string
	// Database to use
	Database string
}

// StorageBackend holds the settings for the used backend component
type StorageBackend struct {
	// Name of the storage backend to use (e.g., mongo)
	Backend string

	MongoBackend MongoBackend
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
	// Specified the maximum age for summoner spells data in minutes until it's invalidated. 0 means it is always fetched newly.
	MaxAgeSummonerSpells uint32
	// Specified the maximum age for items data in minutes until it's invalidated. 0 means it is always fetched newly.
	MaxAgeItems uint32
	// Specifies a default RiotClient for use if not otherwise specified in requests or function calls
	DefaultRiotClient string
}

// FetchRunner holds the settings for the FetchRunner
type FetchRunner struct {
	// Specifies the RiotAPI region to use
	Region string

	// Specified the update interval for fetching Summoner Matches in minutes > 0
	UpdateIntervalSummonerMatches uint32
	// Specified the update interval for fetching Free Rotation in minutes > 0 (disabled if = 0)
	UpdateIntervalFreeRotation uint32

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

	// Specifies if for Summoners encountered in fetched matches an additional fetch run shall be performed (warning, can take a while)
	FetchMatchesForSeenSummoners bool

	// If true stops fetching matches for a summoner if it encounters a game version != latest known game version
	FetchOnlyLatestGameVersion bool

	// Specify what the latest game version for fetching is, see config parameter below for details
	LatestGameVersionForFetching string
}

// ChampionsStats holds the settings for the Champions analysis of the StatsRunner
type ChampionsStats struct {
	Enabled        bool    // Specifies if the ChampionStats calculation shall be activated
	UpdateInverval uint32  // Update Interval for running the SummonerSpellsStats calculations in minutes > 0
	RoleThreshold  float64 // Percent value over which a role is considered relevant for the Champion
}

// ItemsStats holds the settings for the Items analysis of the StatsRunner
type ItemsStats struct {
	Enabled        bool   // Specifies if the ItemsStats calculation shall be activated
	UpdateInverval uint32 // Update Interval for running the ItemsStats calculations in minutes > 0
}

// SummonerSpellsStats holds the settings for the Summoner Spells analysis of the StatsRunner
type SummonerSpellsStats struct {
	Enabled                 bool   // Specifies if the SummonerSpellsStats calculation shall be activated
	UpdateInverval          uint32 // Update Interval for running the SummonerSpellsStats calculations in minutes > 0
	KeepOnlyHighestPickRate bool   // Store only the SummonerSpells combination per role/total with the highest pick rate
}

// RunesReforgedStats holds the settings for the Runes Reforged analysis of the StatsRunner
type RunesReforgedStats struct {
	Enabled                 bool   // Specifies if the RunesReforgedStats calculation shall be activated
	UpdateInverval          uint32 // Update Interval for running the RunesReforgedStats calculations in minutes > 0
	KeepOnlyHighestPickRate bool   // Store only the Runes Reforged combination per role/total with the highest pick rate
}

// StatsRunner holds the settings for the StatsRunner
type StatsRunner struct {
	RunRScripts            bool   // Specifies if R scripts shall be used (needs a running R installation)
	RScriptPath            string // Path to the R scripts (distributed with alolstats)
	RPlotsOutputPath       string // Path where the generated plots shall be stored
	RScriptsUpdateInterval uint32 // Update Interval for running the R scripts in minutes > 0

	GameVersion []string // We want to do stats calculations for the following versions, must be valid game versions, ordered decending, e.g. 9.5, 9.4, ..., see https://ddragon.leagueoflegends.com/api/versions.json, e.g., 9.1.1, 8.24.1

	ChampionsStats      ChampionsStats      // ChampionsStats worker settings
	ItemsStats          ItemsStats          // ItemsStats worker settings
	SummonerSpellsStats SummonerSpellsStats // SummonerSpells worker settings
	RunesReforgedStats  RunesReforgedStats  // Runes Reforged worker settings
}

// Config holds the complete ALolStats config
type Config struct {
	API            API
	RiotClient     map[string]RiotClient
	LoLStorage     LoLStorage
	StorageBackend StorageBackend
	FetchRunner    map[string]FetchRunner
	StatsRunner    StatsRunner
}
