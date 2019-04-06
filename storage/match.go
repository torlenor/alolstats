package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/torlenor/alolstats/riotclient"
)

// getMatch gets a match from storage or riot client based on GameID
func (s *Storage) getMatchFromClient(client riotclient.Client, id uint64) (riotclient.MatchDTO, error) {
	match, err := s.backend.GetMatch(id)
	if err != nil {
		s.log.Warnln(err)
		match, err := client.MatchByID(id)
		if err != nil {
			s.log.Warnln(err)
			return riotclient.MatchDTO{}, err
		}
		s.log.Debugf("Returned Match %d from Riot API", id)
		s.backend.StoreMatch(match)
		return *match, nil
	}
	s.log.Debugf("Returned Match %d from Storage", id)
	return *match, nil
}

// GetMatch gets a match from storage or riot client based on GameID
func (s *Storage) GetMatch(id uint64) (riotclient.MatchDTO, error) {
	return s.getMatchFromClient(s.riotClient, id)
}

// GetRegionalMatch gets a match from storage or riot client based on GameID for a specific region
func (s *Storage) GetRegionalMatch(region string, id uint64) (riotclient.MatchDTO, error) {
	if client, ok := s.riotClients[region]; ok {
		return s.getMatchFromClient(client, id)
	}
	return riotclient.MatchDTO{}, fmt.Errorf("Invalid region specified: %s", region)
}

// fetchAndStoreMatchFromClient gets a match from Riot Client and stores it in storage backend if it doesn't exist, yet
func (s *Storage) fetchAndStoreMatchFromClient(client riotclient.Client, id uint64) (*riotclient.MatchDTO, error) {
	_, err := s.backend.GetMatch(id)
	if err != nil {
		match, err := client.MatchByID(id)
		if err != nil {
			s.log.Warnln(err)
			return nil, err
		}
		s.log.Debugf("Storing Match %d from Riot API in Backend", id)
		s.backend.StoreMatch(match)
		return match, nil
	}
	return nil, nil
}

// FetchAndStoreMatch gets a match from Riot Client and stores it in storage backend if it doesn't exist, yet
func (s *Storage) FetchAndStoreMatch(id uint64) (*riotclient.MatchDTO, error) {
	return s.fetchAndStoreMatchFromClient(s.riotClient, id)
}

// RegionalFetchAndStoreMatch gets a match from Riot Client for a specific region and stores it in storage backend if it doesn't exist, yet
func (s *Storage) RegionalFetchAndStoreMatch(region string, id uint64) (*riotclient.MatchDTO, error) {
	if client, ok := s.riotClients[region]; ok {
		return s.fetchAndStoreMatchFromClient(client, id)
	}
	return nil, fmt.Errorf("Invalid region specified: %s", region)
}

// GetStoredMatchesByGameVersionAndChampionID gets all matches for a specific game version and Champion ID
func (s *Storage) GetStoredMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (riotclient.Matches, error) {
	matches, err := s.backend.GetMatchesByGameVersionAndChampionID(gameVersion, championID)
	return *matches, err
}

// GetStoredMatchesByGameVersionChampionIDMapBetweenQueueIDs gets all matches for a specific game version, Champion ID, map id and gtequeue <= queue id <= ltequeue
func (s *Storage) GetStoredMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (riotclient.Matches, error) {
	matches, err := s.backend.GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion, championID, mapID, ltequeue, gtequeue)
	return *matches, err
}

// GetStoredMatchesCursorByGameVersion returns cursor to matches specific to a certain game version
func (s *Storage) GetStoredMatchesCursorByGameVersion(gameVersion string) (QueryCursor, error) {
	return s.backend.GetMatchesCursorByGameVersion(gameVersion)
}

// GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs returns cursor to matches specific to a certain game version
func (s *Storage) GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (QueryCursor, error) {
	return s.backend.GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs(gameVersion, championID, mapID, ltequeue, gtequeue)
}

// GetMatchesCursorByGameVersionMapBetweenQueueIDs returns cursor to matches specific to a certain game version
func (s *Storage) GetMatchesCursorByGameVersionMapBetweenQueueIDs(gameVersion string, mapID uint64, ltequeue uint64, gtequeue uint64) (QueryCursor, error) {
	return s.backend.GetMatchesCursorByGameVersionMapBetweenQueueIDs(gameVersion, mapID, ltequeue, gtequeue)
}

// GetMatchesCursorByGameVersionMapQueueID returns cursor to matches specific to a certain game version
func (s *Storage) GetMatchesCursorByGameVersionMapQueueID(gameVersion string, mapID uint64, queueid uint64) (QueryCursor, error) {
	return s.backend.GetMatchesCursorByGameVersionMapQueueID(gameVersion, mapID, queueid)
}

// getMatchesByAccountIDFromClient gets all match references for a specified Account ID and startIndex, endIndex
func (s *Storage) getMatchesByAccountIDFromClient(client riotclient.Client, accountID string, beginIndex uint32, endIndex uint32) (*riotclient.MatchlistDTO, error) {
	beginIndexStr := strconv.FormatInt(int64(beginIndex), 10)
	endIndexStr := strconv.FormatInt(int64(endIndex), 10)
	return client.MatchesByAccountID(accountID, map[string]string{"beginIndex": beginIndexStr, "endIndex": endIndexStr})
}

// GetMatchesByAccountID gets all match references for a specified Account ID and startIndex, endIndex
func (s *Storage) GetMatchesByAccountID(accountID string, beginIndex uint32, endIndex uint32) (*riotclient.MatchlistDTO, error) {
	return s.getMatchesByAccountIDFromClient(s.riotClient, accountID, beginIndex, endIndex)
}

// GetRegionalMatchesByAccountID gets all match references for a specified Account ID and startIndex, endIndex for a specific region
func (s *Storage) GetRegionalMatchesByAccountID(region string, accountID string, beginIndex uint32, endIndex uint32) (*riotclient.MatchlistDTO, error) {
	if client, ok := s.riotClients[region]; ok {
		return s.getMatchesByAccountIDFromClient(client, accountID, beginIndex, endIndex)
	}
	return nil, fmt.Errorf("Invalid region specified: %s", region)
}

func (s *Storage) getMatchEndpoint(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln("Received Rest API Match request from", r.RemoteAddr)
	if val, ok := r.URL.Query()["id"]; ok {
		if len(val) == 0 {
			s.log.Warnf("id parameter was empty in request")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			s.log.Warnf("Could not convert value %s to GameID", val)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		match, err := s.GetMatch(id)
		if err != nil {
			s.log.Warnf("Could not get match for id=%d", id)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		out, err := json.Marshal(match)
		if err != nil {
			s.log.Errorln(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		io.WriteString(w, string(out))
	}

	atomic.AddUint64(&s.stats.handledRequests, 1)
}
