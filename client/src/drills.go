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

type MathQuestion struct {
	ID          string
	Topic       string
	Prompt      string
	Options     []string
	AnswerIndex int
	Explanation string
}

var mathBank = []MathQuestion{
	{
		ID:          "quad-reciprocal",
		Topic:       "Quadratic equations",
		Prompt:      "If α and β are the roots of x² - 7x + 10 = 0, what is 1/α + 1/β?",
		Options:     []string{"3/5", "7/10", "10/7", "5/7"},
		AnswerIndex: 1,
		Explanation: "For x² - 7x + 10 = 0: α + β = 7 and αβ = 10. So 1/α + 1/β = (α + β)/(αβ) = 7/10.",
	},
	{
		ID:          "log-balance",
		Topic:       "Logarithms",
		Prompt:      "Solve log₂(x - 1) + log₂(x + 1) = 3.",
		Options:     []string{"x = 2", "x = 3", "x = √7", "x = 4"},
		AnswerIndex: 1,
		Explanation: "Combine logs: log₂((x - 1)(x + 1)) = 3, so x² - 1 = 8, x² = 9, therefore x = 3.",
	},
	{
		ID:          "trig-max",
		Topic:       "Trigonometry",
		Prompt:      "For real θ, the maximum value of sin(θ) + cos(θ) is:",
		Options:     []string{"1", "√2", "2", "3/2"},
		AnswerIndex: 1,
		Explanation: "sin(θ) + cos(θ) = √2 sin(θ + π/4), which attains maximum √2.",
	},
	{
		ID:          "gp-sum",
		Topic:       "Sequences and series",
		Prompt:      "The sum to infinity of a GP is 18 and first term is 6. What is common ratio r?",
		Options:     []string{"1/3", "2/3", "3/4", "1/2"},
		AnswerIndex: 1,
		Explanation: "S = a / (1 - r) => 18 = 6 / (1 - r) => 1 - r = 1/3 => r = 2/3.",
	},
	{
		ID:          "derivative-slope",
		Topic:       "Differentiation",
		Prompt:      "What is the slope of tangent to y = x³ - 3x² + 2x at x = 2?",
		Options:     []string{"2", "4", "6", "8"},
		AnswerIndex: 0,
		Explanation: "dy/dx = 3x² - 6x + 2. At x = 2: 3(4) - 6(2) + 2 = 12 - 12 + 2 = 2.",
	},
	{
		ID:          "definite-integral",
		Topic:       "Integration",
		Prompt:      "Evaluate ∫ (3x² + 2x) dx from x = 0 to 1.",
		Options:     []string{"1", "3/2", "2", "5/2"},
		AnswerIndex: 2,
		Explanation: "[x³ + x²] from 0 to 1 = (1 + 1) - 0 = 2.",
	},
	{
		ID:          "determinant",
		Topic:       "Matrices",
		Prompt:      "What is the determinant of [[2, 3], [5, 7]]?",
		Options:     []string{"-1", "-2", "1", "2"},
		AnswerIndex: 0,
		Explanation: "det = (2 * 7) - (3 * 5) = 14 - 15 = -1.",
	},
	{
		ID:          "complex-product",
		Topic:       "Complex numbers",
		Prompt:      "What is the value of (3 + 4i)(3 - 4i)?",
		Options:     []string{"7", "9", "16", "25"},
		AnswerIndex: 3,
		Explanation: "Use conjugate product a² + b² = 3² + 4² = 9 + 16 = 25.",
	},
}

type MathState struct {
	Running       bool
	Score         int
	Streak        int
	Multiplier    float64
	TimeRemaining int
	CurrentQ      *MathQuestion
	StopTicker    chan struct{}
}

var mState = &MathState{
	Multiplier: 1.0,
}

var mathInitialized = false

func initMathView() {
	if mathInitialized {
		return
	}
	mathInitialized = true

	doc := js.Global().Get("document")

	// Start button listener
	startBtn := doc.Call("getElementById", "mathStartBtn")
	if !startBtn.IsNull() {
		startBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			startMathSprint()
			return nil
		}))
	}
}

func startMathSprint() {
	if mState.Running {
		return
	}

	mState.Running = true
	mState.Score = 0
	mState.Streak = 0
	mState.Multiplier = 1.0
	mState.TimeRemaining = 45
	mState.StopTicker = make(chan struct{})

	updateMathHUD()
	nextMathQuestion()

	// 45s countdown timer in background goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-mState.StopTicker:
				return
			case <-ticker.C:
				mState.TimeRemaining--
				updateMathHUD()
				if mState.TimeRemaining <= 0 {
					finishMathSprint()
					return
				}
			}
		}
	}()
}

func nextMathQuestion() {
	rand.Seed(time.Now().UnixNano())
	q := mathBank[rand.Intn(len(mathBank))]
	mState.CurrentQ = &q

	doc := js.Global().Get("document")
	promptEl := doc.Call("getElementById", "mathEquationPrompt")
	if !promptEl.IsNull() {
		promptEl.Set("innerText", q.Prompt)
		promptEl.Get("style").Set("fontSize", "1.65rem")
	}

	// Render multiple choice option buttons into container
	feedbackEl := doc.Call("getElementById", "mathFeedbackNote")
	if !feedbackEl.IsNull() {
		var sb strings.Builder
		sb.WriteString(`<div style="display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px; width: 100%; max-width: 440px; margin: 12px auto 0;">`)
		for i, opt := range q.Options {
			sb.WriteString(fmt.Sprintf(`
				<button class="btn btn-secondary math-opt-btn" data-opt-idx="%d" style="height: 48px; border-radius: 12px; font-weight: 700; font-size: 1.05rem;">
					%s
				</button>
			`, i, opt))
		}
		sb.WriteString(`</div>`)
		feedbackEl.Set("innerHTML", sb.String())

		// Bind option clicks
		btns := feedbackEl.Call("querySelectorAll", ".math-opt-btn")
		bCount := btns.Get("length").Int()
		for i := 0; i < bCount; i++ {
			btn := btns.Index(i)
			btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				idx := this.Call("getAttribute", "data-opt-idx").Int()
				handleMathAnswer(idx)
				return nil
			}))
		}
	}
}

func handleMathAnswer(chosenIdx int) {
	if !mState.Running || mState.CurrentQ == nil {
		return
	}

	if chosenIdx == mState.CurrentQ.AnswerIndex {
		// Correct
		mState.Streak++
		if mState.Streak >= 10 {
			mState.Multiplier = 3.0
		} else if mState.Streak >= 6 {
			mState.Multiplier = 2.0
		} else if mState.Streak >= 3 {
			mState.Multiplier = 1.5
		}

		pts := int(100.0 * mState.Multiplier)
		mState.Score += pts
		triggerToast("Correct Answer!", fmt.Sprintf("+%d Points (Multiplier: %.1fx)", pts, mState.Multiplier), "success")
	} else {
		// Incorrect
		mState.Streak = 0
		mState.Multiplier = 1.0
		triggerToast("Incorrect", mState.CurrentQ.Explanation, "warning")
	}

	updateMathHUD()
	nextMathQuestion()
}

func updateMathHUD() {
	doc := js.Global().Get("document")
	setElementText(doc, "mathLiveScore", fmt.Sprintf("%d", mState.Score))
	setElementText(doc, "mathLiveTime", fmt.Sprintf("%ds", mState.TimeRemaining))
	setElementText(doc, "mathLiveStreak", fmt.Sprintf("%d 🔥", mState.Streak))
	setElementText(doc, "mathLiveMultiplier", fmt.Sprintf("%.1fx", mState.Multiplier))
}

func finishMathSprint() {
	mState.Running = false
	if mState.StopTicker != nil {
		close(mState.StopTicker)
	}

	doc := js.Global().Get("document")
	promptEl := doc.Call("getElementById", "mathEquationPrompt")
	if !promptEl.IsNull() {
		promptEl.Set("innerText", fmt.Sprintf("Sprint Finished! Final Score: %d", mState.Score))
		promptEl.Get("style").Set("fontSize", "1.8rem")
	}

	feedbackEl := doc.Call("getElementById", "mathFeedbackNote")
	if !feedbackEl.IsNull() {
		feedbackEl.Set("innerHTML", `<div style="font-size: 0.95rem; color: #10b981; font-weight: 700; margin-top: 14px;">Results saved to operator session profile.</div>`)
	}

	// Save result
	res := Result{
		ActivityType: "math",
		Score:        float64(mState.Score),
		Accuracy:     100.0,
		DurationMs:   45000,
		Timestamp:    time.Now().UnixMilli(),
	}
	SaveActivityResult(res)

	triggerToast("Math Sprint Concluded", fmt.Sprintf("Final Score: %d Points", mState.Score), "success")
}
