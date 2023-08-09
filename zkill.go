package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Killmail struct {
	KillmailID int `json:"killmail_id"`
	ZKB        struct {
		Hash string `json:"hash"`
	} `json:"zkb"`
}

type KillmailStats struct {
	Type       string `json:"type"`
	ID         int    `json:"id"`
	TopAllTime []struct {
		Type string `json:"type"`
		Data []struct {
			Kills      int `json:"kills"`
			ShipTypeID int `json:"shipTypeID"`
		} `json:"data"`
	} `json:"topAllTime"`
}

var killmailCache = make(map[string][]Killmail)
var killmailCached = make(map[string]bool)

var killmailStatsCache = make(map[string]KillmailStats)

func GetRecentLosses(characterID int, shipID int) ([]Killmail, error) {
	if killmailCached[fmt.Sprintf("%d_%d", characterID, shipID)] {
		return getRecentLossesFromCache(characterID, shipID), nil
	}

	url := fmt.Sprintf("https://zkillboard.com/api/losses/characterID/%d/shipTypeID/%d/", characterID, shipID)
	killmails, err := fetchRecentLossesFromAPI(url)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%d_%d", characterID, shipID)
	killmailCache[key] = killmails
	killmailCached[key] = true

	if len(killmails) > 10 {
		return killmails[:10], nil
	} else {
		return killmails, nil
	}
}

func getRecentLossesFromCache(characterID int, shipID int) []Killmail {
	if killmailCached[fmt.Sprintf("%d_%d", characterID, shipID)] {
		return killmailCache[fmt.Sprintf("%d_%d", characterID, shipID)]
	}
	return nil
}

func fetchRecentLossesFromAPI(url string) ([]Killmail, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var killmails []Killmail
	err = json.NewDecoder(resp.Body).Decode(&killmails)
	if err != nil {
		return nil, err
	}

	return killmails, nil
}

func GetTopShips(characterID int) ([]int, error) {
	if killmailStats, ok := killmailStatsCache[fmt.Sprintf("%d", characterID)]; ok {
		var result []int
		for _, topShip := range killmailStats.TopAllTime[4].Data {
			result = append(result, topShip.ShipTypeID)
		}
		return result, nil
	}

	url := fmt.Sprintf("https://zkillboard.com/api/stats/characterID/%d/", characterID)
	killmailStats, err := fetchTopShipsFromAPI(url)
	if err != nil {
		return nil, err
	}

	var result []int
	for _, topShip := range killmailStats.TopAllTime[4].Data {
		result = append(result, topShip.ShipTypeID)
	}
	return result, nil
}

func fetchTopShipsFromAPI(url string) (KillmailStats, error) {
	resp, err := http.Get(url)
	if err != nil {
		return KillmailStats{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var killmailStats KillmailStats
	err = json.NewDecoder(resp.Body).Decode(&killmailStats)
	if err != nil {
		return KillmailStats{}, err
	}

	return killmailStats, nil
}
