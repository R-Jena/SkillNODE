//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"
)

const (
	ProfileStorageKey = "skillnode_profile"
	ResultsStorageKey = "skillnode_results"
	ThemeStorageKey   = "skillnode-theme"
)

type Profile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Region    string `json:"region"`
	Bio       string `json:"bio"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type Result struct {
	ID           string  `json:"id"`
	ActivityType string  `json:"activityType"`
	Score        float64 `json:"score"`
	Accuracy     float64 `json:"accuracy"`
	DurationMs   int64   `json:"durationMs"`
	Timestamp    int64   `json:"timestamp"`
}

var currentProfile *Profile

func initStore() {
	currentProfile = loadProfileFromStorage()
	if currentProfile == nil {
		currentProfile = createDefaultGuestProfile()
		saveProfileToStorage(currentProfile)
	}
}

func loadProfileFromStorage() *Profile {
	storage := js.Global().Get("localStorage")
	raw := storage.Call("getItem", ProfileStorageKey)
	if raw.IsNull() || raw.IsUndefined() || raw.String() == "" {
		return nil
	}

	var p Profile
	err := json.Unmarshal([]byte(raw.String()), &p)
	if err != nil || p.ID == "" {
		return nil
	}
	return &p
}

func saveProfileToStorage(p *Profile) {
	if p == nil {
		return
	}
	currentProfile = p
	data, err := json.Marshal(p)
	if err != nil {
		return
	}
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", ProfileStorageKey, string(data))
}

func createDefaultGuestProfile() *Profile {
	now := time.Now().UnixMilli()
	return &Profile{
		ID:        fmt.Sprintf("player_%d", now),
		Name:      "Operator",
		Username:  "operator",
		Email:     "operator@skillnode.dev",
		Region:    "Global Edge Node",
		Bio:       "Local-first skill operator. Focused on high-speed execution, algorithmic problem solving, and zero-latency competition.",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func GetSessionProfile() *Profile {
	if currentProfile == nil {
		initStore()
	}
	return currentProfile
}

func LoadResults() []Result {
	storage := js.Global().Get("localStorage")
	raw := storage.Call("getItem", ResultsStorageKey)
	if raw.IsNull() || raw.IsUndefined() || raw.String() == "" {
		return []Result{}
	}

	var list []Result
	_ = json.Unmarshal([]byte(raw.String()), &list)
	return list
}

func SaveActivityResult(r Result) {
	if r.ID == "" {
		r.ID = fmt.Sprintf("res_%d", time.Now().UnixMilli())
	}
	if r.Timestamp == 0 {
		r.Timestamp = time.Now().UnixMilli()
	}

	existing := LoadResults()
	// Keep newest at front, cap at 50 results
	updated := append([]Result{r}, existing...)
	if len(updated) > 50 {
		updated = updated[:50]
	}

	data, err := json.Marshal(updated)
	if err == nil {
		storage := js.Global().Get("localStorage")
		storage.Call("setItem", ResultsStorageKey, string(data))
	}
}

func GetStoredTheme() string {
	storage := js.Global().Get("localStorage")
	t := storage.Call("getItem", ThemeStorageKey)
	if t.IsNull() || t.IsUndefined() || t.String() == "" {
		return "dark-purple"
	}
	return t.String()
}

func SetStoredTheme(theme string) {
	storage := js.Global().Get("localStorage")
	storage.Call("setItem", ThemeStorageKey, theme)
}
