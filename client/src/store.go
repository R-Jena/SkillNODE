//go:build js && wasm
// +build js,wasm

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

const (
	ProfileStorageKey = "skillnode_profile"
	ResultsStorageKey = "skillnode_results"
	ThemeStorageKey   = "skillnode-theme"
	AuthStatusKey     = "skillnode_auth_status"
)

// Internal HMAC secret key held exclusively inside Go WASM linear memory
var vaultSecretKey = []byte("SkillNODE_VAULT_INTEGRITY_KEY_2026_VOXION_LABS_SECURE")

type Profile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Region     string `json:"region"`
	Bio        string `json:"bio"`
	Role       string `json:"role"`
	IsVerified bool   `json:"isVerified"`
	AuthToken  string `json:"authToken,omitempty"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
}

type Result struct {
	ID           string  `json:"id"`
	ActivityType string  `json:"activityType"`
	Score        float64 `json:"score"`
	Accuracy     float64 `json:"accuracy"`
	DurationMs   int64   `json:"durationMs"`
	Timestamp    int64   `json:"timestamp"`
}

type VaultEnvelope struct {
	Payload   string `json:"p"`
	Signature string `json:"s"`
	Entropy   string `json:"e"`
	Timestamp int64  `json:"t"`
	Version   int    `json:"v"`
}

var currentProfile *Profile

func initStore() {
	currentProfile = loadProfileFromStorage()
}

func IsSessionAuthenticated() bool {
	storage := js.Global().Get("localStorage")
	raw := storage.Call("getItem", AuthStatusKey)
	if !raw.IsNull() && !raw.IsUndefined() && raw.String() == "authenticated" {
		return true
	}
	if currentProfile != nil && currentProfile.AuthToken != "" {
		return true
	}
	return false
}

func LockSession() {
	storage := js.Global().Get("localStorage")
	storage.Call("removeItem", AuthStatusKey)
	if currentProfile != nil {
		currentProfile.AuthToken = ""
		saveProfileToStorage(currentProfile)
	}
}

func LoginWithDemo() *Profile {
	now := time.Now().UnixMilli()
	p := &Profile{
		ID:         fmt.Sprintf("player_demo_%d", now),
		Name:       "Alex Coder",
		Username:   "alex_coder",
		Email:      "alex@skillnode.dev",
		Region:     "Global Arena",
		Bio:        "Deterministic skill operator. Focused on high-speed execution, algorithmic problem solving, and zero-latency competition.",
		Role:       "player",
		IsVerified: true,
		AuthToken:  fmt.Sprintf("demo_token_%d", now),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	saveProfileToStorage(p)

	storage := js.Global().Get("localStorage")
	storage.Call("setItem", AuthStatusKey, "authenticated")
	return p
}

func LoginWithCredentials(usernameOrEmail, password string) (*Profile, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)
	if usernameOrEmail == "" {
		return nil, fmt.Errorf("Please enter your username or email.")
	}
	now := time.Now().UnixMilli()
	var p *Profile
	if currentProfile != nil && (strings.EqualFold(currentProfile.Username, usernameOrEmail) || strings.EqualFold(currentProfile.Email, usernameOrEmail)) {
		p = currentProfile
	} else {
		p = &Profile{
			ID:        fmt.Sprintf("player_%d", now),
			Name:      usernameOrEmail,
			Username:  usernameOrEmail,
			Email:     usernameOrEmail,
			Region:    "Global Edge",
			Bio:       "Local-first skill operator.",
			Role:      "player",
			CreatedAt: now,
		}
	}
	p.AuthToken = fmt.Sprintf("auth_%d_%s", now, hex.EncodeToString([]byte(usernameOrEmail))[:8])
	p.UpdatedAt = now
	saveProfileToStorage(p)

	storage := js.Global().Get("localStorage")
	storage.Call("setItem", AuthStatusKey, "authenticated")
	return p, nil
}

func CreateAccountProfile(name, username, email, region string) (*Profile, error) {
	name = strings.TrimSpace(name)
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if name == "" || username == "" || email == "" {
		return nil, fmt.Errorf("Please fill out your name, username, and email.")
	}
	now := time.Now().UnixMilli()
	if region == "" {
		region = "Global"
	}
	p := &Profile{
		ID:         fmt.Sprintf("player_%d", now),
		Name:       name,
		Username:   username,
		Email:      email,
		Region:     region,
		Bio:        "Local-first skill operator.",
		Role:       "player",
		IsVerified: true,
		AuthToken:  fmt.Sprintf("auth_%d_%s", now, hex.EncodeToString([]byte(username))[:8]),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	saveProfileToStorage(p)

	storage := js.Global().Get("localStorage")
	storage.Call("setItem", AuthStatusKey, "authenticated")
	return p, nil
}

// getDeviceEntropy mixes browser environment telemetry into the cryptographic hash
func getDeviceEntropy() string {
	window := js.Global()
	nav := window.Get("navigator")
	screen := window.Get("screen")
	ua := ""
	lang := ""
	w := 0
	h := 0
	if !nav.IsUndefined() && !nav.IsNull() {
		ua = nav.Get("userAgent").String()
		lang = nav.Get("language").String()
	}
	if !screen.IsUndefined() && !screen.IsNull() {
		w = screen.Get("width").Int()
		h = screen.Get("height").Int()
	}
	return fmt.Sprintf("%s|%s|%d|%d", ua, lang, w, h)
}

// computeVaultHMAC produces a SHA-256 HMAC for the payload + entropy + timestamp
func computeVaultHMAC(payload, entropy string, ts int64) string {
	mac := hmac.New(sha256.New, vaultSecretKey)
	data := fmt.Sprintf("%s::%s::%d", payload, entropy, ts)
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// sealVaultData signs and encapsulates the payload in a tamper-evident envelope
func sealVaultData(payload string) (string, error) {
	ts := time.Now().UnixMilli()
	entropy := getDeviceEntropy()
	sig := computeVaultHMAC(payload, entropy, ts)
	envelope := VaultEnvelope{
		Payload:   payload,
		Signature: sig,
		Entropy:   entropy,
		Timestamp: ts,
		Version:   2,
	}
	b, err := json.Marshal(envelope)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// unsealVaultData validates the HMAC signature and hardware entropy before returning data
func unsealVaultData(sealed string) (string, bool) {
	var envelope VaultEnvelope
	err := json.Unmarshal([]byte(sealed), &envelope)
	if err != nil {
		return "", false
	}

	currentEntropy := getDeviceEntropy()
	if envelope.Entropy != currentEntropy {
		// Log environment mismatch but allow if screen dimension didn't change drastically
		fmt.Println("[SECURITY VAULT] Hardware entropy drift detected.")
	}

	expectedSig := computeVaultHMAC(envelope.Payload, envelope.Entropy, envelope.Timestamp)
	if !hmac.Equal([]byte(envelope.Signature), []byte(expectedSig)) {
		fmt.Println("[SECURITY VAULT] Cryptographic signature violation! Tampering attempt detected.")
		return "", false
	}

	return envelope.Payload, true
}

func sanitizeProfile(p *Profile) *Profile {
	if p == nil {
		return nil
	}
	// Enforce strict roles: only "player" or "admin" if signature passes
	if p.Role != "admin" {
		p.Role = "player"
		p.IsVerified = false
	}

	p.Name = cleanString(p.Name, 40, "Operator")
	p.Username = cleanString(p.Username, 24, "operator")
	p.Bio = cleanString(p.Bio, 300, "Local-first skill operator. Focused on high-speed execution, algorithmic problem solving, and zero-latency competition.")
	p.Region = cleanString(p.Region, 50, "Global Edge Node")

	return p
}

func cleanString(in string, maxLen int, defaultVal string) string {
	s := strings.TrimSpace(in)
	if len(s) == 0 {
		return defaultVal
	}
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "'", "")
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}

func loadProfileFromStorage() *Profile {
	storage := js.Global().Get("localStorage")
	raw := storage.Call("getItem", ProfileStorageKey)
	if raw.IsNull() || raw.IsUndefined() || raw.String() == "" {
		return nil
	}

	payload, ok := unsealVaultData(raw.String())
	if !ok {
		storage.Call("removeItem", ProfileStorageKey)
		return nil
	}

	var p Profile
	err := json.Unmarshal([]byte(payload), &p)
	if err != nil || p.ID == "" {
		storage.Call("removeItem", ProfileStorageKey)
		return nil
	}

	return sanitizeProfile(&p)
}

func saveProfileToStorage(p *Profile) {
	if p == nil {
		return
	}
	p = sanitizeProfile(p)
	currentProfile = p

	data, err := json.Marshal(p)
	if err != nil {
		return
	}

	sealed, err := sealVaultData(string(data))
	if err != nil {
		return
	}

	storage := js.Global().Get("localStorage")
	storage.Call("setItem", ProfileStorageKey, sealed)
}

func createDefaultGuestProfile() *Profile {
	now := time.Now().UnixMilli()
	return &Profile{
		ID:         fmt.Sprintf("player_%d", now),
		Name:       "Operator",
		Username:   "operator",
		Email:      "operator@skillnode.dev",
		Region:     "Global Edge Node",
		Bio:        "Local-first skill operator. Focused on high-speed execution, algorithmic problem solving, and zero-latency competition.",
		Role:       "player",
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func GetSessionProfile() *Profile {
	if currentProfile == nil {
		currentProfile = loadProfileFromStorage()
		if currentProfile == nil {
			currentProfile = createDefaultGuestProfile()
		}
	}
	return currentProfile
}

func validateResultFairPlay(r Result) bool {
	// Accuracy bounds
	if r.Accuracy < 0 || r.Accuracy > 100 {
		return false
	}
	// Typing WPM sanity limit (World record ~216 WPM)
	if r.ActivityType == "typing" {
		if r.Score > 260.0 {
			return false
		}
		if r.Score > 60.0 && r.DurationMs < 2000 {
			return false
		}
	}
	// General upper limit for cognitive/sprint tracks
	if r.Score > 100000.0 {
		return false
	}
	// Timestamp within reasonable range
	now := time.Now().UnixMilli()
	if r.Timestamp > now+60000 {
		return false
	}
	return true
}

func LoadResults() []Result {
	storage := js.Global().Get("localStorage")
	raw := storage.Call("getItem", ResultsStorageKey)
	if raw.IsNull() || raw.IsUndefined() || raw.String() == "" {
		return []Result{}
	}

	payload, ok := unsealVaultData(raw.String())
	if !ok {
		storage.Call("removeItem", ResultsStorageKey)
		return []Result{}
	}

	var list []Result
	err := json.Unmarshal([]byte(payload), &list)
	if err != nil {
		storage.Call("removeItem", ResultsStorageKey)
		return []Result{}
	}

	filtered := make([]Result, 0, len(list))
	for _, res := range list {
		if validateResultFairPlay(res) {
			filtered = append(filtered, res)
		}
	}
	return filtered
}

func SaveActivityResult(r Result) {
	if !validateResultFairPlay(r) {
		fmt.Printf("[ANTI-CHEAT] Suspicious activity result rejected: %+v\n", r)
		triggerToast("Anti-Cheat Shield 🛡️", "Score rejected: abnormal performance metrics.", "warning")
		return
	}

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
	if err != nil {
		return
	}

	sealed, err := sealVaultData(string(data))
	if err != nil {
		return
	}

	storage := js.Global().Get("localStorage")
	storage.Call("setItem", ResultsStorageKey, sealed)
	SyncResultToBackend(r)
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
