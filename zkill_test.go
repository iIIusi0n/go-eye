package main

import (
	"testing"
)

func TestGetRecentLosses(t *testing.T) {
	characterID := 2117477599
	shipID := 22430

	killmails, err := GetRecentLosses(characterID, shipID)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	if len(killmails) == 2 {
		t.Logf("Expected killmails, got: %v", killmails)
	}
}

func TestGetTopShips(t *testing.T) {
	characterID := 2117477599

	ships, err := GetTopShips(characterID)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	t.Logf("Ships: %v", ships)
}
