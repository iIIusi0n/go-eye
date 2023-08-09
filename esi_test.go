package main

import (
	"reflect"
	"testing"
)

func TestResolveIdsToNames(t *testing.T) {
	ids := []int{29990, 602}
	expectedNames := []string{"Loki", "Kestrel"}

	names, err := ResolveIdsToNames(ids)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf("Expected names: %v, got: %v", expectedNames, names)
	}

	t.Logf("Names: %v", names)
}

func TestResolveNamesToCharacterIDs(t *testing.T) {
	names := []string{"Market Scammer", "Market Trickster"}
	expectedIDs := []int{2117477599, 2118503862}

	ids, err := ResolveNamesToCharacterIDs(names)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	if !reflect.DeepEqual(ids, expectedIDs) {
		t.Errorf("Expected IDs: %v, got: %v", expectedIDs, ids)
	}

	t.Logf("IDs: %v", ids)
}

func TestGetItemsFromKillmail(t *testing.T) {
	killmailID := 97332126
	killmailHash := "1627401de883aa21f99c8618e5b8ca59f7904dae"
	expectedItems := []int{16650, 3828, 11399, 2393, 32880, 3683, 1319, 1319, 16644, 16648, 16273, 1319, 16647, 16649, 29984, 16274, 16638, 16641, 16651}

	items, _, err := GetItemsFromKillmail(killmailID, killmailHash)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Expected items: %v, got: %v", expectedItems, items)
	}

	t.Logf("Items: %v", items)
}

func TestGetItemsFromKillmailCaching(t *testing.T) {
	killmailID := 97332126
	killmailHash := "1627401de883aa21f99c8618e5b8ca59f7904dae"
	expectedItems := []int{16650, 3828, 11399, 2393, 32880, 3683, 1319, 1319, 16644, 16648, 16273, 1319, 16647, 16649, 29984, 16274, 16638, 16641, 16651}

	items, _, err := GetItemsFromKillmail(killmailID, killmailHash)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Expected items: %v, got: %v", expectedItems, items)
	}

	t.Logf("Items: %v", items)
}

func TestResolveItemNamesToIDs(t *testing.T) {
	names := []string{"Loki", "Kestrel"}
	expectedIDs := []int{29990, 602}

	ids, err := ResolveItemNamesToIDs(names)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	if !reflect.DeepEqual(ids, expectedIDs) {
		t.Errorf("Expected IDs: %v, got: %v", expectedIDs, ids)
	}

	t.Logf("IDs: %v", ids)
}
