//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

type FermiQuestion struct {
	ID          string
	Prompt      string
	Unit        string
	TrueValue   float64
	Context     string
	Explanation string
}

var fermiBank = []FermiQuestion{
	{
		ID:          "seconds-year",
		Prompt:      "How many seconds are in a standard 365-day year?",
		Unit:        "seconds",
		TrueValue:   31536000,
		Context:     "365 days × 24 hours × 60 minutes × 60 seconds.",
		Explanation: "Calculation: 365 × 86,400 = 31,536,000 seconds.",
	},
	{
		ID:          "speed-sound",
		Prompt:      "What is the speed of sound in dry air at 20°C in meters per second?",
		Unit:        "m/s",
		TrueValue:   343,
		Context:     "Mach 1 at standard sea-level atmospheric conditions.",
		Explanation: "At 20°C (68°F), the speed of sound in air is approximately 343 m/s (1,235 km/h).",
	},
	{
		ID:          "earth-moon-km",
		Prompt:      "What is the distance from Earth to the Moon at perigee (closest point) in kilometers?",
		Unit:        "km",
		TrueValue:   363300,
		Context:     "Lunar orbit varies from 363,300 km (perigee) to 405,500 km (apogee).",
		Explanation: "The perigee distance is roughly 363,300 km, giving an average distance of ~384,400 km.",
	},
	{
		ID:          "piano-keys",
		Prompt:      "How many keys are on a standard modern acoustic grand piano?",
		Unit:        "keys",
		TrueValue:   88,
		Context:     "52 white keys (naturals) and 36 black keys (accidentals).",
		Explanation: "Modern pianos standardized on 88 keys (7 octaves plus a minor third from A0 to C8).",
	},
	{
		ID:          "human-bones",
		Prompt:      "How many bones are in the typical adult human skeleton?",
		Unit:        "bones",
		TrueValue:   206,
		Context:     "Infants are born with ~270 bones, which fuse as adulthood is reached.",
		Explanation: "An adult human has exactly 206 bones after cranial and pelvic fusion.",
	},
	{
		ID:          "olympic-pool-liters",
		Prompt:      "What is the approximate water volume of an Olympic swimming pool in liters?",
		Unit:        "liters",
		TrueValue:   2500000,
		Context:     "Dimensions: 50m long × 25m wide × 2m nominal depth.",
		Explanation: "50m × 25m × 2m = 2,500 cubic meters = 2,500,000 liters (approx 660,000 gallons).",
	},
	{
		ID:          "periodic-elements",
		Prompt:      "How many confirmed chemical elements are recognized on the periodic table?",
		Unit:        "elements",
		TrueValue:   118,
		Context:     "From Hydrogen (1) through Oganesson (118).",
		Explanation: "The periodic table features 118 confirmed elements, completed with IUPAC additions in 2016.",
	},
	{
		ID:          "marathon-meters",
		Prompt:      "What is the official total distance of a standard marathon in meters?",
		Unit:        "meters",
		TrueValue:   42195,
		Context:     "26 miles and 385 yards, standardized at the 1908 London Olympics.",
		Explanation: "The standard distance is 42.195 km = 42,195 meters.",
	},
	{
		ID:          "speed-light-kms",
		Prompt:      "What is the speed of light in a vacuum in kilometers per second (rounded)?",
		Unit:        "km/s",
		TrueValue:   299792,
		Context:     "Universal physical constant c = 299,792,458 m/s.",
		Explanation: "In kilometers per second, c ≈ 299,792 km/s (~300,000 km/s).",
	},
	{
		ID:          "leap-year-minutes",
		Prompt:      "How many total minutes are in one full leap year (366 days)?",
		Unit:        "minutes",
		TrueValue:   527040,
		Context:     "366 days × 24 hours/day × 60 minutes/hour.",
		Explanation: "Calculation: 366 × 1,440 = 527,040 minutes.",
	},
	{
		ID:          "ethernet-mtu",
		Prompt:      "What is the standard Ethernet Maximum Transmission Unit (MTU) payload size in bytes?",
		Unit:        "bytes",
		TrueValue:   1500,
		Context:     "Standard IEEE 802.3 Ethernet frame data payload limit.",
		Explanation: "Standard Ethernet MTU is 1,500 bytes (excluding 14-byte Ethernet header and 4-byte CRC).",
	},
	{
		ID:          "sha256-bits",
		Prompt:      "How many bits comprise a standard SHA-256 cryptographic digest output?",
		Unit:        "bits",
		TrueValue:   256,
		Context:     "Secure Hash Algorithm 256-bit fixed output length.",
		Explanation: "SHA-256 generates a 256-bit (32-byte) digest, rendered as a 64-character hexadecimal string.",
	},
	{
		ID:          "octopus-hearts",
		Prompt:      "How many hearts does an octopus possess?",
		Unit:        "hearts",
		TrueValue:   3,
		Context:     "Two branchial hearts pump blood through gills; one systemic heart pumps to the body.",
		Explanation: "Octopuses have 3 hearts. Two pump blood through each of the two gills, and one circulates blood to the organs.",
	},
	{
		ID:          "tcp-port-count",
		Prompt:      "What is the total number of distinct TCP/UDP port numbers available in IPv4/IPv6?",
		Unit:        "ports",
		TrueValue:   65536,
		Context:     "Ports are defined by a 16-bit integer field (0 to 65535).",
		Explanation: "2^16 = 65,536 total available ports (0 to 65,535).",
	},
	{
		ID:          "earth-equator-km",
		Prompt:      "What is the circumference of the Earth at the equator in kilometers (rounded)?",
		Unit:        "km",
		TrueValue:   40075,
		Context:     "Equatorial radius is ~6,378 km.",
		Explanation: "2 × π × 6,378 km ≈ 40,075 km around the equator.",
	},
	{
		ID:          "pcie4-x16-gbps",
		Prompt:      "What is the theoretical maximum bidirectional bandwidth of PCIe 4.0 x16 in GB/s?",
		Unit:        "GB/s",
		TrueValue:   64,
		Context:     "PCIe 4.0 delivers ~2 GB/s per lane across 16 bidirectional lanes.",
		Explanation: "16 lanes × 1.969 GB/s per lane ≈ 31.5 GB/s unidirectional, ~64 GB/s full duplex.",
	},
	{
		ID:          "hours-100-days",
		Prompt:      "How many hours are in exactly 100 days?",
		Unit:        "hours",
		TrueValue:   2400,
		Context:     "100 days × 24 hours per day.",
		Explanation: "100 × 24 = 2,400 hours.",
	},
	{
		ID:          "ipv4-total-addresses",
		Prompt:      "How many total theoretical IPv4 addresses exist in billions (approx 2^32)?",
		Unit:        "billion",
		TrueValue:   4.29,
		Context:     "A 32-bit address space yields 4,294,967,296 distinct addresses.",
		Explanation: "2^32 = 4,294,967,296 addresses ≈ 4.29 billion addresses.",
	},
}

type CrypticWord struct {
	Word string
	Clue string
}

var crypticBank = []CrypticWord{
	{Word: "MUTEX", Clue: "Mutual exclusion synchronization primitive used to prevent data races."},
	{Word: "CIPHER", Clue: "Algorithm for performing encryption or decryption on plaintext."},
	{Word: "KERNEL", Clue: "The core computer program managing system resources and hardware interaction."},
	{Word: "THREAD", Clue: "Smallest sequence of programmed instructions managed independently by a scheduler."},
	{Word: "VECTOR", Clue: "Dynamic array sequence container capable of amortized O(1) appending."},
	{Word: "ROUTER", Clue: "Networking device forwarding data packets between distinct IP networks."},
	{Word: "SOCKET", Clue: "Endpoint for sending or receiving network data across a communication channel."},
	{Word: "BUFFER", Clue: "Temporary memory storage holding data while it is transferred between locations."},
	{Word: "PRISMA", Clue: "Next-generation ORM and query builder for deterministic database interaction."},
	{Word: "SCHEMA", Clue: "Structural definition specifying tables, columns, constraints, and relationships."},
}

type GuessState struct {
	ActiveTab string // "fermi", "binary", "word"

	// Fermi Mode
	CurrentFermiIdx int
	FermiScore      int
	FermiQuestionsDone int

	// Binary Search Mode
	BinaryTarget  int
	BinaryMin     int
	BinaryMax     int
	BinaryGuesses int
	BinarySolved  bool
	BinaryLog     []string

	// Cryptic Word Mode
	CurrentWordIdx int
	WordGuesses    int
	WordSolved     bool
	WordLog        []string
}

var gState = &GuessState{
	ActiveTab: "fermi",
}
var guessInitialized = false

func initGuessView() {
	if guessInitialized {
		return
	}
	guessInitialized = true

	doc := js.Global().Get("document")

	// Tab switcher
	tabs := doc.Call("querySelectorAll", "[data-guess-tab]")
	tCount := tabs.Get("length").Int()
	for i := 0; i < tCount; i++ {
		tab := tabs.Index(i)
		tab.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			mode := this.Call("getAttribute", "data-guess-tab").String()
			switchGuessTab(mode)
			return nil
		}))
	}

	// Fermi submit button
	fermiBtn := doc.Call("getElementById", "fermiSubmitBtn")
	if !fermiBtn.IsNull() {
		fermiBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			submitFermiGuess()
			return nil
		}))
	}

	// Fermi next button
	fermiNextBtn := doc.Call("getElementById", "fermiNextBtn")
	if !fermiNextBtn.IsNull() {
		fermiNextBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			nextFermiQuestion()
			return nil
		}))
	}

	// Binary search submit button
	binBtn := doc.Call("getElementById", "binarySubmitBtn")
	if !binBtn.IsNull() {
		binBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			submitBinaryGuess()
			return nil
		}))
	}

	// Binary search reset button
	binResetBtn := doc.Call("getElementById", "binaryResetBtn")
	if !binResetBtn.IsNull() {
		binResetBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resetBinaryGame()
			return nil
		}))
	}

	// Word submit button
	wordBtn := doc.Call("getElementById", "wordSubmitBtn")
	if !wordBtn.IsNull() {
		wordBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			submitWordGuess()
			return nil
		}))
	}

	// Word next button
	wordResetBtn := doc.Call("getElementById", "wordResetBtn")
	if !wordResetBtn.IsNull() {
		wordResetBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resetWordGame()
			return nil
		}))
	}

	// Initialize Fermi Mode
	resetFermiGame()
	resetBinaryGame()
	resetWordGame()
	switchGuessTab("fermi")
}

func switchGuessTab(mode string) {
	gState.ActiveTab = mode
	doc := js.Global().Get("document")

	tabs := doc.Call("querySelectorAll", "[data-guess-tab]")
	tCount := tabs.Get("length").Int()
	for i := 0; i < tCount; i++ {
		t := tabs.Index(i)
		if t.Call("getAttribute", "data-guess-tab").String() == mode {
			t.Get("classList").Call("add", "active")
		} else {
			t.Get("classList").Call("remove", "active")
		}
	}

	setPanelDisplay(doc, "panel-guess-fermi", mode == "fermi")
	setPanelDisplay(doc, "panel-guess-binary", mode == "binary")
	setPanelDisplay(doc, "panel-guess-word", mode == "word")
}

func setPanelDisplay(doc js.Value, id string, show bool) {
	el := doc.Call("getElementById", id)
	if !el.IsNull() {
		if show {
			el.Get("style").Set("display", "flex")
		} else {
			el.Get("style").Set("display", "none")
		}
	}
}

// ── Fermi Mode Logic ──

func resetFermiGame() {
	gState.FermiScore = 0
	gState.FermiQuestionsDone = 0
	gState.CurrentFermiIdx = 0
	rand.Seed(time.Now().UnixNano())
	renderCurrentFermi()
}

func renderCurrentFermi() {
	if len(fermiBank) == 0 {
		return
	}
	doc := js.Global().Get("document")
	q := fermiBank[gState.CurrentFermiIdx%len(fermiBank)]

	setElementText(doc, "fermiPrompt", q.Prompt)
	setElementText(doc, "fermiUnitLabel", fmt.Sprintf("Units: %s", q.Unit))
	setElementText(doc, "fermiContext", q.Context)
	setElementText(doc, "fermiLiveScore", fmt.Sprintf("%d", gState.FermiScore))
	setElementText(doc, "fermiProgress", fmt.Sprintf("Question %d of %d", (gState.CurrentFermiIdx%len(fermiBank))+1, len(fermiBank)))

	inputEl := doc.Call("getElementById", "fermiInput")
	if !inputEl.IsNull() {
		inputEl.Set("value", "")
		inputEl.Get("style").Set("display", "block")
	}

	btnSubmit := doc.Call("getElementById", "fermiSubmitBtn")
	if !btnSubmit.IsNull() {
		btnSubmit.Get("style").Set("display", "inline-flex")
	}

	btnNext := doc.Call("getElementById", "fermiNextBtn")
	if !btnNext.IsNull() {
		btnNext.Get("style").Set("display", "none")
	}

	feedbackEl := doc.Call("getElementById", "fermiFeedback")
	if !feedbackEl.IsNull() {
		feedbackEl.Set("innerHTML", `<div style="font-size: 0.85rem; color: var(--text-faint);">Enter your numerical estimate and test your heuristic accuracy.</div>`)
	}
}

func submitFermiGuess() {
	doc := js.Global().Get("document")
	inputEl := doc.Call("getElementById", "fermiInput")
	if inputEl.IsNull() {
		return
	}
	valStr := strings.TrimSpace(inputEl.Get("value").String())
	guessVal, err := strconv.ParseFloat(valStr, 64)
	if err != nil || guessVal <= 0 {
		triggerToast("Invalid Input", "Please enter a positive numeric value.", "warning")
		return
	}

	q := fermiBank[gState.CurrentFermiIdx%len(fermiBank)]
	errorPct := math.Abs(guessVal-q.TrueValue) / q.TrueValue * 100.0

	pts := 0
	rating := "Order of magnitude discrepancy"
	badgeColor := "#ef4444"

	if errorPct <= 5.0 {
		pts = 150
		rating = "Bullseye Accuracy! 🎯"
		badgeColor = "#10b981"
	} else if errorPct <= 20.0 {
		pts = 100
		rating = "Exceptional Precision ⚡"
		badgeColor = "#34d399"
	} else if errorPct <= 50.0 {
		pts = 60
		rating = "Solid Engineering Heuristic 🔍"
		badgeColor = "#fbbf24"
	} else if errorPct <= 100.0 {
		pts = 30
		rating = "Reasonable Ballpark 📊"
		badgeColor = "#f59e0b"
	} else {
		pts = 10
		rating = "Heuristic Divergence"
		badgeColor = "#f472b6"
	}

	gState.FermiScore += pts
	gState.FermiQuestionsDone++

	setElementText(doc, "fermiLiveScore", fmt.Sprintf("%d", gState.FermiScore))

	feedbackEl := doc.Call("getElementById", "fermiFeedback")
	if !feedbackEl.IsNull() {
		html := fmt.Sprintf(`
			<div style="padding: 16px; border-radius: 12px; background: var(--surface-2); border: 1px solid var(--border-subtle); display: flex; flex-direction: column; gap: 8px;">
				<div style="display: flex; align-items: center; justify-content: space-between;">
					<strong style="color: %s; font-size: 0.95rem;">%s (+%d Points)</strong>
					<span style="font-size: 0.8rem; color: var(--text-faint);">Error: %.1f%%</span>
				</div>
				<div style="font-size: 0.85rem; color: var(--text); line-height: 1.5;">
					<strong>True Value:</strong> %s %s | <strong>Your Guess:</strong> %s %s
				</div>
				<div style="font-size: 0.8rem; color: var(--text-muted); line-height: 1.5; margin-top: 4px;">
					%s
				</div>
			</div>
		`, badgeColor, rating, pts, errorPct, formatNumber(q.TrueValue), q.Unit, formatNumber(guessVal), q.Unit, q.Explanation)
		feedbackEl.Set("innerHTML", html)
	}

	btnSubmit := doc.Call("getElementById", "fermiSubmitBtn")
	if !btnSubmit.IsNull() {
		btnSubmit.Get("style").Set("display", "none")
	}

	btnNext := doc.Call("getElementById", "fermiNextBtn")
	if !btnNext.IsNull() {
		btnNext.Get("style").Set("display", "inline-flex")
	}

	// Save progress
	SaveActivityResult(Result{
		ActivityType: "guess",
		Score:        float64(gState.FermiScore),
		Accuracy:     math.Max(0, 100.0-errorPct),
		DurationMs:   15000,
		Timestamp:    time.Now().UnixMilli(),
	})
	triggerToast("Heuristic Evaluated", fmt.Sprintf("+%d Points (Error: %.1f%%)", pts, errorPct), "success")
}

func nextFermiQuestion() {
	gState.CurrentFermiIdx = (gState.CurrentFermiIdx + 1) % len(fermiBank)
	renderCurrentFermi()
}

// ── Binary Search Mode Logic ──

func resetBinaryGame() {
	rand.Seed(time.Now().UnixNano())
	gState.BinaryTarget = rand.Intn(1000) + 1
	gState.BinaryMin = 1
	gState.BinaryMax = 1000
	gState.BinaryGuesses = 0
	gState.BinarySolved = false
	gState.BinaryLog = []string{}

	doc := js.Global().Get("document")
	setElementText(doc, "binTargetRange", "Range: 1 – 1,000")
	setElementText(doc, "binGuessCount", "0 guesses")
	setElementText(doc, "binSearchStatus", "Enter a number between 1 and 1,000 to begin binary search deduction.")

	inputEl := doc.Call("getElementById", "binInput")
	if !inputEl.IsNull() {
		inputEl.Set("value", "")
		inputEl.Set("disabled", false)
	}

	btnSubmit := doc.Call("getElementById", "binarySubmitBtn")
	if !btnSubmit.IsNull() {
		btnSubmit.Set("disabled", false)
	}

	renderBinaryLog()
}

func submitBinaryGuess() {
	if gState.BinarySolved {
		return
	}
	doc := js.Global().Get("document")
	inputEl := doc.Call("getElementById", "binInput")
	if inputEl.IsNull() {
		return
	}
	valStr := strings.TrimSpace(inputEl.Get("value").String())
	guess, err := strconv.Atoi(valStr)
	if err != nil || guess < 1 || guess > 1000 {
		triggerToast("Out of Bounds", "Please enter an integer between 1 and 1,000.", "warning")
		return
	}

	gState.BinaryGuesses++
	diff := int(math.Abs(float64(guess - gState.BinaryTarget)))

	var hint string
	var statusMsg string
	if guess == gState.BinaryTarget {
		gState.BinarySolved = true
		pts := int(math.Max(1000.0-float64(gState.BinaryGuesses-10)*80.0, 150.0))
		statusMsg = fmt.Sprintf("🎉 Target %d Acquired in %d guesses! (Theoretical optimal: 10 steps). Score: %d pts.", gState.BinaryTarget, gState.BinaryGuesses, pts)
		hint = "🎯 MATCH"

		inputEl.Set("disabled", true)
		btnSubmit := doc.Call("getElementById", "binarySubmitBtn")
		if !btnSubmit.IsNull() {
			btnSubmit.Set("disabled", true)
		}

		SaveActivityResult(Result{
			ActivityType: "guess",
			Score:        float64(pts),
			Accuracy:     math.Max(0, 100.0-float64(gState.BinaryGuesses-10)*5.0),
			DurationMs:   int64(gState.BinaryGuesses * 2500),
			Timestamp:    time.Now().UnixMilli(),
		})
		triggerToast("Target Acquired! 🎯", fmt.Sprintf("Found %d in %d guesses (+%d pts)", gState.BinaryTarget, gState.BinaryGuesses, pts), "success")
	} else if guess < gState.BinaryTarget {
		if guess >= gState.BinaryMin {
			gState.BinaryMin = guess + 1
		}
		hotness := getHotness(diff)
		hint = fmt.Sprintf("⬆️ HIGHER (%s)", hotness)
		statusMsg = fmt.Sprintf("Target is HIGHER than %d. Active Search Interval: [%d, %d]", guess, gState.BinaryMin, gState.BinaryMax)
	} else {
		if guess <= gState.BinaryMax {
			gState.BinaryMax = guess - 1
		}
		hotness := getHotness(diff)
		hint = fmt.Sprintf("⬇️ LOWER (%s)", hotness)
		statusMsg = fmt.Sprintf("Target is LOWER than %d. Active Search Interval: [%d, %d]", guess, gState.BinaryMin, gState.BinaryMax)
	}

	gState.BinaryLog = append([]string{fmt.Sprintf("Guess #%d: %d — %s", gState.BinaryGuesses, guess, hint)}, gState.BinaryLog...)

	setElementText(doc, "binTargetRange", fmt.Sprintf("Search Space: [%d – %d]", gState.BinaryMin, gState.BinaryMax))
	setElementText(doc, "binGuessCount", fmt.Sprintf("%d guesses (log₂ limit: 10)", gState.BinaryGuesses))
	setElementText(doc, "binSearchStatus", statusMsg)

	inputEl.Set("value", "")
	renderBinaryLog()
}

func getHotness(diff int) string {
	if diff <= 15 {
		return "Burning Hot 🔥"
	} else if diff <= 50 {
		return "Very Warm ⚡"
	} else if diff <= 150 {
		return "Warm"
	} else if diff <= 300 {
		return "Cool"
	}
	return "Freezing Cold ❄️"
}

func renderBinaryLog() {
	doc := js.Global().Get("document")
	logEl := doc.Call("getElementById", "binLogFeed")
	if logEl.IsNull() {
		return
	}
	if len(gState.BinaryLog) == 0 {
		logEl.Set("innerHTML", `<div style="font-size: 0.8rem; color: var(--text-faint); padding: 8px; text-align: center;">No guesses registered yet.</div>`)
		return
	}

	var sb strings.Builder
	for _, entry := range gState.BinaryLog {
		color := "var(--text-muted)"
		if strings.Contains(entry, "MATCH") {
			color = "#10b981"
		} else if strings.Contains(entry, "Burning") {
			color = "#f472b6"
		}
		sb.WriteString(fmt.Sprintf(`
			<div style="padding: 6px 12px; border-radius: 8px; background: var(--surface-2); font-size: 0.82rem; font-family: monospace; color: %s; display: flex; justify-content: space-between;">
				<span>%s</span>
			</div>
		`, color, entry))
	}
	logEl.Set("innerHTML", sb.String())
}

// ── Cryptic Word Deduction Mode Logic ──

func resetWordGame() {
	rand.Seed(time.Now().UnixNano())
	gState.CurrentWordIdx = rand.Intn(len(crypticBank))
	gState.WordGuesses = 0
	gState.WordSolved = false
	gState.WordLog = []string{}

	doc := js.Global().Get("document")
	w := crypticBank[gState.CurrentWordIdx]
	setElementText(doc, "wordCluePrompt", w.Clue)
	setElementText(doc, "wordLengthHint", fmt.Sprintf("Length: %d characters", len(w.Word)))

	inputEl := doc.Call("getElementById", "wordInput")
	if !inputEl.IsNull() {
		inputEl.Set("value", "")
		inputEl.Set("disabled", false)
	}

	btnSubmit := doc.Call("getElementById", "wordSubmitBtn")
	if !btnSubmit.IsNull() {
		btnSubmit.Set("disabled", false)
	}

	renderWordLog()
}

func submitWordGuess() {
	if gState.WordSolved {
		return
	}
	doc := js.Global().Get("document")
	inputEl := doc.Call("getElementById", "wordInput")
	if inputEl.IsNull() {
		return
	}

	guess := strings.ToUpper(strings.TrimSpace(inputEl.Get("value").String()))
	target := crypticBank[gState.CurrentWordIdx].Word

	if len(guess) != len(target) {
		triggerToast("Length Mismatch", fmt.Sprintf("Target word has %d letters. You entered %d.", len(target), len(guess)), "warning")
		return
	}

	gState.WordGuesses++

	// Build letter tiles feedback
	var tileSb strings.Builder
	tileSb.WriteString(`<div style="display: flex; gap: 6px; margin: 4px 0;">`)
	allCorrect := true

	for i := 0; i < len(guess); i++ {
		ch := string(guess[i])
		targetCh := string(target[i])
		bg := "rgba(255, 255, 255, 0.08)"
		fg := "var(--text-muted)"

		if ch == targetCh {
			bg = "#10b981"
			fg = "#ffffff"
		} else if strings.Contains(target, ch) {
			bg = "#f59e0b"
			fg = "#ffffff"
			allCorrect = false
		} else {
			allCorrect = false
		}

		tileSb.WriteString(fmt.Sprintf(`
			<span style="display: inline-flex; width: 32px; height: 32px; align-items: center; justify-content: center; font-weight: 800; border-radius: 6px; background: %s; color: %s; font-family: monospace;">
				%s
			</span>
		`, bg, fg, ch))
	}
	tileSb.WriteString(`</div>`)

	gState.WordLog = append([]string{tileSb.String()}, gState.WordLog...)

	if allCorrect {
		gState.WordSolved = true
		pts := int(math.Max(500.0-float64(gState.WordGuesses-1)*50.0, 100.0))
		inputEl.Set("disabled", true)
		btnSubmit := doc.Call("getElementById", "wordSubmitBtn")
		if !btnSubmit.IsNull() {
			btnSubmit.Set("disabled", true)
		}

		SaveActivityResult(Result{
			ActivityType: "guess",
			Score:        float64(pts),
			Accuracy:     100.0,
			DurationMs:   int64(gState.WordGuesses * 4000),
			Timestamp:    time.Now().UnixMilli(),
		})
		triggerToast("Word Solved! 🎉", fmt.Sprintf("Identified %s in %d attempts (+%d pts)", target, gState.WordGuesses, pts), "success")
	}

	inputEl.Set("value", "")
	renderWordLog()
}

func renderWordLog() {
	doc := js.Global().Get("document")
	logEl := doc.Call("getElementById", "wordLogFeed")
	if logEl.IsNull() {
		return
	}
	if len(gState.WordLog) == 0 {
		logEl.Set("innerHTML", `<div style="font-size: 0.8rem; color: var(--text-faint); padding: 8px; text-align: center;">Enter your first attempt above.</div>`)
		return
	}

	var sb strings.Builder
	for _, entry := range gState.WordLog {
		sb.WriteString(entry)
	}
	logEl.Set("innerHTML", sb.String())
}

func formatNumber(n float64) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.2fM", n/1000000.0)
	} else if n >= 1000 {
		return fmt.Sprintf("%.0f", n)
	}
	return fmt.Sprintf("%.2f", n)
}
