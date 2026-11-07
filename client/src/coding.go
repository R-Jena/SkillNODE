//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
	"time"
)

type TestCase struct {
	ArgsJSON string
	Expected string
}

type CodeChallenge struct {
	ID           string
	Title        string
	Difficulty   string
	Description  string
	FunctionName string
	Language     string
	Starter      string
	Constraints  []string
	Tests        []TestCase
}

var codeChallenges = []CodeChallenge{
	{
		ID:           "mergeIntervals",
		Title:        "Merge Busy Intervals",
		Difficulty:   "Advanced",
		Language:     "javascript",
		FunctionName: "mergeIntervals",
		Description:  "Merge overlapping booking intervals and return them in sorted order.",
		Starter: `function mergeIntervals(intervals) {
  if (!intervals.length) return [];
  intervals.sort((a, b) => a[0] - b[0]);
  const merged = [intervals[0]];
  for (let i = 1; i < intervals.length; i++) {
    const prev = merged[merged.length - 1];
    const curr = intervals[i];
    if (curr[0] <= prev[1]) {
      prev[1] = Math.max(prev[1], curr[1]);
    } else {
      merged.push(curr);
    }
  }
  return merged;
}`,
		Constraints: []string{"Sort before merging.", "Time target: O(n log n)."},
		Tests: []TestCase{
			{ArgsJSON: "[[[1,3],[2,6],[8,10]]]", Expected: "[[1,6],[8,10]]"},
			{ArgsJSON: "[[[1,4],[4,5]]]", Expected: "[[1,5]]"},
			{ArgsJSON: "[[[5,7],[1,2],[2,4]]]", Expected: "[[1,4],[5,7]]"},
		},
	},
	{
		ID:           "matrixDiagonalGap",
		Title:        "Matrix Diagonal Gap",
		Difficulty:   "Medium",
		Language:     "javascript",
		FunctionName: "matrixDiagonalGap",
		Description:  "Return the absolute difference between the primary and secondary diagonal sums.",
		Starter: `function matrixDiagonalGap(grid) {
  let primary = 0, secondary = 0;
  const n = grid.length;
  for (let i = 0; i < n; i++) {
    primary += grid[i][i];
    secondary += grid[i][n - 1 - i];
  }
  return Math.abs(primary - secondary);
}`,
		Constraints: []string{"A single pass is sufficient: O(n)."},
		Tests: []TestCase{
			{ArgsJSON: "[[[[1,2,3],[4,5,6],[7,8,9]]]]", Expected: "0"},
			{ArgsJSON: "[[[[6,1,1],[4,5,2],[9,3,7]]]]", Expected: "4"},
		},
	},
	{
		ID:           "longestPeak",
		Title:        "Longest Peak Length",
		Difficulty:   "Advanced",
		Language:     "python",
		FunctionName: "",
		Description:  "Return the length of the longest peak in an integer array.",
		Starter: `def longest_peak(arr):
    # A peak strictly rises and strictly falls
    longest = 0
    i = 1
    while i < len(arr) - 1:
        is_peak = arr[i-1] < arr[i] and arr[i] > arr[i+1]
        if not is_peak:
            i += 1
            continue
        left = i - 2
        while left >= 0 and arr[left] < arr[left+1]:
            left -= 1
        right = i + 2
        while right < len(arr) and arr[right] < arr[right-1]:
            right += 1
        length = right - left - 1
        if length > longest:
            longest = length
        i = right
    return longest`,
		Constraints: []string{"Time complexity target: O(n).", "Space complexity: O(1)."},
		Tests: []TestCase{
			{ArgsJSON: "review-mode", Expected: "Checkpoint: Strict peak detection arr[i-1] < arr[i] > arr[i+1]"},
			{ArgsJSON: "review-mode", Expected: "Checkpoint: Expand left and right outward from peak vertex"},
		},
	},
}

var activeChallenge = &codeChallenges[0]
var codingInitialized = false

func initCodingView() {
	if codingInitialized {
		renderActiveChallenge()
		return
	}
	codingInitialized = true

	doc := js.Global().Get("document")

	// 1. Language selector tabs
	tabs := doc.Call("querySelectorAll", "[data-lang]")
	tCount := tabs.Get("length").Int()
	for i := 0; i < tCount; i++ {
		tab := tabs.Index(i)
		tab.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			lang := this.Call("getAttribute", "data-lang").String()
			selectLanguage(lang)
			return nil
		}))
	}

	// 2. Run Tests button
	runBtn := doc.Call("getElementById", "codeRunBtn")
	if !runBtn.IsNull() {
		runBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			executeCodeTests()
			return nil
		}))
	}

	// 3. Reset Starter button
	resetBtn := doc.Call("getElementById", "codeResetStarterBtn")
	if !resetBtn.IsNull() {
		resetBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			editor := doc.Call("getElementById", "codeEditorTextarea")
			if !editor.IsNull() && activeChallenge != nil {
				editor.Set("value", activeChallenge.Starter)
			}
			return nil
		}))
	}

	renderActiveChallenge()
}

func selectLanguage(lang string) {
	doc := js.Global().Get("document")

	// Update active tab style
	tabs := doc.Call("querySelectorAll", "[data-lang]")
	tCount := tabs.Get("length").Int()
	for i := 0; i < tCount; i++ {
		t := tabs.Index(i)
		if t.Call("getAttribute", "data-lang").String() == lang {
			t.Get("classList").Call("add", "active")
		} else {
			t.Get("classList").Call("remove", "active")
		}
	}

	// Find first challenge for language
	for i := range codeChallenges {
		if codeChallenges[i].Language == lang {
			activeChallenge = &codeChallenges[i]
			renderActiveChallenge()
			return
		}
	}

	// Default to first challenge if none for language
	activeChallenge = &codeChallenges[0]
	renderActiveChallenge()
}

func renderActiveChallenge() {
	if activeChallenge == nil {
		return
	}

	doc := js.Global().Get("document")

	setElementText(doc, "challengeTitle", activeChallenge.Title)
	setElementText(doc, "challengeDescription", activeChallenge.Description)
	setElementText(doc, "challengeDifficulty", activeChallenge.Difficulty)
	setElementText(doc, "editorLanguageLabel", fmt.Sprintf("solution.%s", getFileExtension(activeChallenge.Language)))

	// Constraints
	constraintsEl := doc.Call("getElementById", "challengeConstraints")
	if !constraintsEl.IsNull() {
		var sb strings.Builder
		for _, c := range activeChallenge.Constraints {
			sb.WriteString(fmt.Sprintf("<li>%s</li>", c))
		}
		constraintsEl.Set("innerHTML", sb.String())
	}

	// Starter code
	editor := doc.Call("getElementById", "codeEditorTextarea")
	if !editor.IsNull() {
		editor.Set("value", activeChallenge.Starter)
	}

	// Reset test results panel
	resultsFeed := doc.Call("getElementById", "codeTestResultsFeed")
	if !resultsFeed.IsNull() {
		resultsFeed.Set("innerHTML", `<div style="font-size: 0.8rem; color: var(--text-faint); padding: 10px; background: var(--surface-2); border-radius: 8px; text-align: center;">Click 'Run Tests' to validate your solution against sample assertions.</div>`)
	}
}

func executeCodeTests() {
	doc := js.Global().Get("document")
	editor := doc.Call("getElementById", "codeEditorTextarea")
	if editor.IsNull() || activeChallenge == nil {
		return
	}

	userCode := editor.Get("value").String()
	resultsFeed := doc.Call("getElementById", "codeTestResultsFeed")
	if resultsFeed.IsNull() {
		return
	}

	// If JavaScript, evaluate directly via browser JS engine
	if activeChallenge.Language == "javascript" && activeChallenge.FunctionName != "" {
		passCount := 0
		totalTests := len(activeChallenge.Tests)
		var sb strings.Builder

		evalFuncScript := fmt.Sprintf(`
			try {
				%s
				if (typeof %s !== 'function') {
					return { error: 'Expected function named %s' };
				}
				const fn = %s;
				const tests = %s;
				const results = tests.map((t, idx) => {
					try {
						const actual = fn(...JSON.parse(t.ArgsJSON));
						const expected = JSON.parse(t.Expected);
						const pass = JSON.stringify(actual) === JSON.stringify(expected);
						return { pass, actual: JSON.stringify(actual), expected: t.Expected };
					} catch(e) {
						return { pass: false, error: e.message };
					}
				});
				return { results };
			} catch(e) {
				return { error: e.message };
			}
		`, userCode, activeChallenge.FunctionName, activeChallenge.FunctionName, activeChallenge.FunctionName, marshalTests(activeChallenge.Tests))

		evaluator := js.Global().Get("Function").New(evalFuncScript)
		resVal := evaluator.Invoke()

		if !resVal.Get("error").IsUndefined() {
			errStr := resVal.Get("error").String()
			sb.WriteString(fmt.Sprintf(`
				<div style="padding: 10px 14px; border-radius: 10px; background: rgba(239, 68, 68, 0.15); border: 1px solid rgba(239, 68, 68, 0.3); color: #fca5a5; font-size: 0.82rem;">
					<strong>Runtime Error:</strong> %s
				</div>
			`, errStr))
			resultsFeed.Set("innerHTML", sb.String())
			triggerToast("Code Execution Error", errStr, "error")
			return
		}

		resultsArr := resVal.Get("results")
		rLen := resultsArr.Get("length").Int()

		for i := 0; i < rLen; i++ {
			item := resultsArr.Index(i)
			passed := item.Get("pass").Bool()

			if passed {
				passCount++
				sb.WriteString(fmt.Sprintf(`
					<div style="display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; border-radius: 10px; background: rgba(16, 185, 129, 0.12); border: 1px solid rgba(16, 185, 129, 0.25); font-size: 0.82rem;">
						<div style="display: flex; align-items: center; gap: 8px;">
							<i class="fas fa-check-circle" style="color: #10b981;"></i>
							<strong style="color: var(--text);">Test Case %d</strong>
						</div>
						<span style="color: #10b981; font-weight: 700;">PASSED</span>
					</div>
				`, i+1))
			} else {
				expected := item.Get("expected").String()
				actual := item.Get("actual").String()
				sb.WriteString(fmt.Sprintf(`
					<div style="padding: 10px 14px; border-radius: 10px; background: rgba(239, 68, 68, 0.12); border: 1px solid rgba(239, 68, 68, 0.25); font-size: 0.82rem;">
						<div style="display: flex; align-items: center; justify-content: space-between;">
							<strong style="color: #ef4444;">Test Case %d Failed</strong>
							<span style="color: #ef4444; font-weight: 700;">FAIL</span>
						</div>
						<div style="margin-top: 4px; font-size: 0.75rem; color: var(--text-muted);">
							Expected: %s | Got: %s
						</div>
					</div>
				`, i+1, expected, actual))
			}
		}

		resultsFeed.Set("innerHTML", sb.String())

		// Record result in store
		score := (float64(passCount) / float64(totalTests)) * 100.0
		SaveActivityResult(Result{
			ActivityType: "coding",
			Score:        score,
			Accuracy:     score,
			DurationMs:   10000,
			Timestamp:    time.Now().UnixMilli(),
		})

		if passCount == totalTests {
			triggerToast("All Tests Passed! 🎉", fmt.Sprintf("%d/%d test cases verified.", passCount, totalTests), "success")
		} else {
			triggerToast("Tests Completed", fmt.Sprintf("%d/%d test cases passed.", passCount, totalTests), "warning")
		}
	} else {
		// Guided checklist for Python/C++/TS
		var sb strings.Builder
		for i, t := range activeChallenge.Tests {
			sb.WriteString(fmt.Sprintf(`
				<div style="padding: 10px 14px; border-radius: 10px; background: var(--surface-2); border: 1px solid var(--border-subtle); font-size: 0.82rem;">
					<div style="display: flex; align-items: center; gap: 8px;">
						<i class="fas fa-check" style="color: var(--accent-bright);"></i>
						<strong style="color: var(--text);">Review %d:</strong>
						<span style="color: var(--text-muted);">%s</span>
					</div>
				</div>
			`, i+1, t.Expected))
		}
		resultsFeed.Set("innerHTML", sb.String())
		triggerToast("Evaluation Checklist Ready", "Review your solution against key targets.", "info")
	}
}

func marshalTests(tests []TestCase) string {
	data, _ := json.Marshal(tests)
	return string(data)
}

func getFileExtension(lang string) string {
	switch lang {
	case "javascript":
		return "js"
	case "python":
		return "py"
	case "cpp":
		return "cpp"
	case "typescript":
		return "ts"
	case "sql":
		return "sql"
	default:
		return "txt"
	}
}
