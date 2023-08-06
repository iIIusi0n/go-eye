package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// idCache stores the mapping of IDs to names for caching purposes.
var idCache = make(map[int]string)

// nameCache stores the mapping of names to IDs for caching purposes.
var nameCache = make(map[string]int)

// killmailItemsCache stores the items from killmails in cache to avoid repeated API requests.
var killmailItemsCache = make(map[int][]int)

// ResolveIdsToNames resolves a list of IDs to their corresponding names using the cache and EVE Online API.
func ResolveIdsToNames(ids []int) ([]string, error) {
	names, unresolvedIds, err := getNamesFromCache(ids)
	if err != nil {
		return nil, err
	}

	if len(unresolvedIds) > 0 {
		newNames, err := resolveNamesFromAPI(unresolvedIds)
		if err != nil {
			return nil, err
		}
		names = append(names, newNames...)
	}

	return names, nil
}

// getNamesFromCache retrieves names from the cache and returns unresolved IDs.
func getNamesFromCache(ids []int) ([]string, []int, error) {
	names := make([]string, 0)
	unresolvedIds := make([]int, 0)

	for _, id := range ids {
		if name, ok := idCache[id]; ok {
			names = append(names, name)
		} else {
			unresolvedIds = append(unresolvedIds, id)
		}
	}

	return names, unresolvedIds, nil
}

// resolveNamesFromAPI resolves a list of unresolved IDs to names using EVE Online API.
func resolveNamesFromAPI(ids []int) ([]string, error) {
	body, err := json.Marshal(ids)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal unresolved IDs: %w", err)
	}

	req, err := http.NewRequest("POST", "https://esi.evetech.net/latest/universe/names/?datasource=tranquility", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", KUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	dataA, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	var data []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err = json.Unmarshal(dataA, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	newNames := make([]string, 0)
	for _, entry := range data {
		idCache[entry.ID] = entry.Name
		nameCache[entry.Name] = entry.ID
		newNames = append(newNames, entry.Name)
	}

	return newNames, nil
}

// characterInfo holds character ID and name information.
type characterInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ResolveNamesToCharacterIDs resolves a list of names to their corresponding character IDs using the cache and EVE Online API.
func ResolveNamesToCharacterIDs(names []string) ([]int, error) {
	ids, unresolvedNames, err := getIDsFromCache(names)
	if err != nil {
		return nil, err
	}

	if len(unresolvedNames) > 0 {
		newIDs, err := resolveIDsFromAPI(unresolvedNames)
		if err != nil {
			return nil, err
		}
		ids = append(ids, newIDs...)
	}

	return ids, nil
}

// getIDsFromCache retrieves character IDs from the cache and returns unresolved names.
func getIDsFromCache(names []string) ([]int, []string, error) {
	ids := make([]int, 0)
	unresolvedNames := make([]string, 0)

	for _, name := range names {
		if id, ok := nameCache[name]; ok {
			ids = append(ids, id)
		} else {
			unresolvedNames = append(unresolvedNames, name)
		}
	}

	return ids, unresolvedNames, nil
}

// resolveIDsFromAPI resolves a list of unresolved names to character IDs using EVE Online API.
func resolveIDsFromAPI(names []string) ([]int, error) {
	body, err := json.Marshal(names)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal unresolved names: %w", err)
	}

	req, err := http.NewRequest("POST", "https://esi.evetech.net/latest/universe/ids/?datasource=tranquility&language=en", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", KUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	dataA, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	var response struct {
		Characters []characterInfo `json:"characters"`
	}

	err = json.Unmarshal(dataA, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	newIDs := make([]int, 0)
	for _, entry := range response.Characters {
		idCache[entry.ID] = entry.Name
		nameCache[entry.Name] = entry.ID
		newIDs = append(newIDs, entry.ID)
	}

	return newIDs, nil
}

// fetchItemsFromAPI makes an API request and retrieves items from a killmail.
func fetchItemsFromAPI(id int, hash string) ([]int, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://esi.evetech.net/latest/killmails/%d/%s/?datasource=tranquility", id, hash), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", KUserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dataA, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Victim struct {
			Items []struct {
				ItemType int `json:"item_type_id"`
			} `json:"items"`
		} `json:"victim"`
	}

	err = json.Unmarshal(dataA, &data)
	if err != nil {
		return nil, err
	}

	items := make([]int, 0)
	for _, item := range data.Victim.Items {
		items = append(items, item.ItemType)
	}

	return items, nil
}

// GetItemsFromKillmail retrieves items from a killmail with caching support.
func GetItemsFromKillmail(id int, hash string) ([]int, error) {
	// Check if the data is already in the cache
	if cachedItems, ok := killmailItemsCache[id]; ok {
		return cachedItems, nil
	}

	// Fetch items from the API
	items, err := fetchItemsFromAPI(id, hash)
	if err != nil {
		return nil, err
	}

	// Cache the data for future use
	killmailItemsCache[id] = items

	return items, nil
}
