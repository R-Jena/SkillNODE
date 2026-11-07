//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"syscall/js"
	"time"
)

var typingPassages = []string{
	"SkillNODE is where focus becomes visible performance under pressure.",
	"Type fast, think faster, and keep your rhythm clean for the full run.",
	"Sharp hands and calm breathing beat chaotic speed every single time.",
	"Consistency over speed - every keystroke should feel intentional and precise.",
	"The best typists build muscle memory through deliberate, measured practice.",
	"Algorithmic precision and high cadence define mastery in modern competitive arenas.",
	"Maintain cognitive flow and let rhythmic typing lead your technical execution.",
}

var typingPlaceholders = []string{
	"Start typing here. WPM climbs as you keep a steady rhythm.",
	"Type the passage smoothly. Accuracy matters more than panic speed.",
	"Write here to track live speed, accuracy, and overall progress.",
	"Begin your run here and aim for clean, consistent keystrokes.",
}

type TypingState struct {
	Passage       string
	Input         string
	Running       bool
	StartTime     int64
	LastElapsedMs int64
	WPM           float64
	Accuracy      float64
	Progress      float64
	AnomalyAlerts int
	TickerStop    chan struct{}
}

var tState = &TypingState{
	Passage:       typingPassages[0],
	Accuracy:      100,
	AnomalyAlerts: 0,
}

var typingInitialized = false

func initTypingView() {
	if typingInitialized {
		resetTypingRun()
		return
	}
	typingInitialized = true

	doc := js.Global().Get("document")

	// 1. Initial passage render
	resetTypingRun()

	// 2. Input listener
	inputBox := doc.Call("getElementById", "typingInputBox")
	if !inputBox.IsNull() {
		inputBox.Call("addEventListener", "input", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			val := inputBox.Get("value").String()
			handleTypingInput(val)
			return nil
		}))
	}

	// 3. Reset button listener
	resetBtn := doc.Call("getElementById", "typingResetBtn")
	if !resetBtn.IsNull() {
		resetBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			pickNextPassage()
			resetTypingRun()
			return nil
		}))
	}

	// 4. Fair-play blur/focus anomaly listener
	js.Global().Call("addEventListener", "blur", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if tState.Running {
			tState.AnomalyAlerts++
			badge := doc.Call("getElementById", "typingFairPlayBadge")
			if !badge.IsNull() {
				badge.Set("innerHTML", fmt.Sprintf(`<i class="fas fa-exclamation-triangle" style="color: #f59e0b;"></i> Tab Blur (%d)`, tState.AnomalyAlerts))
				badge.Get("style").Set("color", "#f59e0b")
			}
		}
		return nil
	}))

	js.Global().Call("addEventListener", "focus", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		badge := doc.Call("getElementById", "typingFairPlayBadge")
		if !badge.IsNull() && tState.AnomalyAlerts == 0 {
			badge.Set("innerHTML", `<i class="fas fa-shield-alt"></i> Focus Active`)
			badge.Get("style").Set("color", "#10b981")
		}
		return nil
	}))
}

func pickNextPassage() {
	rand.Seed(time.Now().UnixNano())
	next := typingPassages[rand.Intn(len(typingPassages))]
	tState.Passage = next
}

func resetTypingRun() {
	doc := js.Global().Get("document")

	tState.Input = ""
	tState.Running = false
	tState.StartTime = 0
	tState.LastElapsedMs = 0
	tState.WPM = 0
	tState.Accuracy = 100
	tState.Progress = 0

	inputBox := doc.Call("getElementById", "typingInputBox")
	if !inputBox.IsNull() {
		inputBox.Set("value", "")
		randPlaceholder := typingPlaceholders[rand.Intn(len(typingPlaceholders))]
		inputBox.Call("setAttribute", "placeholder", randPlaceholder)
	}

	renderPassageDiff()
	updateTypingHUD(0, 100, 0, 0)
}

func handleTypingInput(val string) {
	if !tState.Running && len(val) > 0 {
		tState.Running = true
		tState.StartTime = time.Now().UnixMilli()
	}

	tState.Input = val

	elapsedMs := int64(0)
	if tState.Running && tState.StartTime > 0 {
		elapsedMs = time.Now().UnixMilli() - tState.StartTime
	}
	tState.LastElapsedMs = elapsedMs

	// Calculate accuracy and metrics
	passage := tState.Passage
	typedLen := len(val)
	passageLen := len(passage)

	correctChars := 0
	for i := 0; i < typedLen && i < passageLen; i++ {
		if val[i] == passage[i] {
			correctChars++
		}
	}

	accuracy := 100.0
	if typedLen > 0 {
		accuracy = (float64(correctChars) / float64(typedLen)) * 100.0
	}

	minutes := float64(elapsedMs) / 60000.0
	words := float64(len(strings.Fields(val)))
	wpm := 0.0
	if minutes > 0 {
		wpm = words / minutes
	}

	progress := 0.0
	if passageLen > 0 {
		progress = (float64(typedLen) / float64(passageLen)) * 100.0
		if progress > 100 {
			progress = 100
		}
	}

	tState.WPM = wpm
	tState.Accuracy = accuracy
	tState.Progress = progress

	// Render diff and update HUD
	renderPassageDiff()
	updateTypingHUD(wpm, accuracy, progress, elapsedMs)

	// Check for completion
	if typedLen >= passageLen && correctChars >= int(float64(passageLen)*0.9) {
		finishTypingRun()
	}
}

func renderPassageDiff() {
	doc := js.Global().Get("document")
	container := doc.Call("getElementById", "typingPassageContainer")
	if container.IsNull() {
		return
	}

	passage := tState.Passage
	input := tState.Input
	inputLen := len(input)

	var sb strings.Builder
	for i, ch := range passage {
		charStr := string(ch)
		if charStr == " " {
			charStr = "&nbsp;"
		}

		if i < inputLen {
			if input[i] == byte(ch) {
				sb.WriteString(fmt.Sprintf(`<span style="color: #10b981; font-weight: 700;">%s</span>`, charStr))
			} else {
				sb.WriteString(fmt.Sprintf(`<span style="color: #ef4444; background: rgba(239, 68, 68, 0.2); text-decoration: underline; font-weight: 700;">%s</span>`, charStr))
			}
		} else if i == inputLen {
			sb.WriteString(fmt.Sprintf(`<span style="background: var(--accent); color: #fff; padding: 1px 2px; border-radius: 3px; font-weight: 700;">%s</span>`, charStr))
		} else {
			sb.WriteString(fmt.Sprintf(`<span style="color: var(--text-faint);">%s</span>`, charStr))
		}
	}

	container.Set("innerHTML", sb.String())
}

func updateTypingHUD(wpm, accuracy, progress float64, elapsedMs int64) {
	doc := js.Global().Get("document")

	setElementText(doc, "typingLiveWpm", fmt.Sprintf("%.0f", wpm))
	setElementText(doc, "typingLiveAccuracy", fmt.Sprintf("%.0f%%", accuracy))
	setElementText(doc, "typingLiveProgress", fmt.Sprintf("%.0f%%", progress))

	seconds := float64(elapsedMs) / 1000.0
	setElementText(doc, "typingLiveElapsed", fmt.Sprintf("%.1fs", seconds))

	progressBar := doc.Call("getElementById", "typingProgressBar")
	if !progressBar.IsNull() {
		progressBar.Get("style").Set("width", fmt.Sprintf("%.1f%%", progress))
	}
}

func finishTypingRun() {
	tState.Running = false

	// Save result
	res := Result{
		ActivityType: "typing",
		Score:        tState.WPM,
		Accuracy:     tState.Accuracy,
		DurationMs:   tState.LastElapsedMs,
		Timestamp:    time.Now().UnixMilli(),
	}
	SaveActivityResult(res)

	// Trigger UI toast
	triggerToast("Typing Run Completed!", fmt.Sprintf("Velocity: %.0f WPM • Accuracy: %.0f%%", tState.WPM, tState.Accuracy), "success")

	// Prepend to session history list
	renderTypingSessionHistory()

	// Pick next passage for immediate flow
	pickNextPassage()
	resetTypingRun()
}

func renderTypingSessionHistory() {
	doc := js.Global().Get("document")
	container := doc.Call("getElementById", "typingRunsList")
	if container.IsNull() {
		return
	}

	results := LoadResults()
	var typingRuns []Result
	for _, r := range results {
		if r.ActivityType == "typing" {
			typingRuns = append(typingRuns, r)
		}
	}

	if len(typingRuns) == 0 {
		return
	}

	var sb strings.Builder
	limit := len(typingRuns)
	if limit > 8 {
		limit = 8
	}

	for i := 0; i < limit; i++ {
		r := typingRuns[i]
		sb.WriteString(fmt.Sprintf(`
			<div style="display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; border-radius: 10px; background: var(--surface-2); border: 1px solid var(--border-subtle); font-size: 0.85rem;">
				<div style="display: flex; align-items: center; gap: 10px;">
					<span style="font-weight: 800; color: var(--accent-bright);">%.0f WPM</span>
					<span style="color: #10b981; font-weight: 600;">%.0f%% Acc</span>
				</div>
				<span style="color: var(--text-faint); font-size: 0.75rem;">%.1fs</span>
			</div>
		`, r.Score, r.Accuracy, float64(r.DurationMs)/1000.0))
	}

	container.Set("innerHTML", sb.String())
}
