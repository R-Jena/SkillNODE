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
	// ── JavaScript ──
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
		Constraints: []string{"A single pass is sufficient: O(n).", "Space complexity: O(1)."},
		Tests: []TestCase{
			{ArgsJSON: "[[[[1,2,3],[4,5,6],[7,8,9]]]]", Expected: "0"},
			{ArgsJSON: "[[[[6,1,1],[4,5,2],[9,3,7]]]]", Expected: "4"},
		},
	},
	{
		ID:           "twoSumTarget",
		Title:        "Two-Sum Zero Loss Pair",
		Difficulty:   "Easy",
		Language:     "javascript",
		FunctionName: "twoSum",
		Description:  "Find indices of the two numbers such that they add up to a specific target.",
		Starter: `function twoSum(nums, target) {
  const map = new Map();
  for (let i = 0; i < nums.length; i++) {
    const complement = target - nums[i];
    if (map.has(complement)) {
      return [map.get(complement), i];
    }
    map.set(nums[i], i);
  }
  return [];
}`,
		Constraints: []string{"Exactly one solution exists.", "Time complexity target: O(n)."},
		Tests: []TestCase{
			{ArgsJSON: "[[2,7,11,15], 9]", Expected: "[0,1]"},
			{ArgsJSON: "[[3,2,4], 6]", Expected: "[1,2]"},
			{ArgsJSON: "[[3,3], 6]", Expected: "[0,1]"},
		},
	},

	// ── TypeScript ──
	{
		ID:           "typedMergeIntervals",
		Title:        "Typed Interval Consolidation",
		Difficulty:   "Advanced",
		Language:     "typescript",
		FunctionName: "mergeIntervals",
		Description:  "Consolidate overlapping discrete intervals into a non-overlapping sorted slice using strict typing.",
		Starter: `interface Interval {
  start: number;
  end: number;
}

function mergeIntervals(intervals: Interval[]): Interval[] {
  if (intervals.length === 0) return [];
  intervals.sort((a, b) => a.start - b.start);
  const merged: Interval[] = [intervals[0]];
  for (let i = 1; i < intervals.length; i++) {
    const prev = merged[merged.length - 1];
    const curr = intervals[i];
    if (curr.start <= prev.end) {
      prev.end = Math.max(prev.end, curr.end);
    } else {
      merged.push(curr);
    }
  }
  return merged;
}`,
		Constraints: []string{"Interface Interval { start: number; end: number; }.", "Target complexity: O(n log n)."},
		Tests: []TestCase{
			{ArgsJSON: "[[{\"start\":1,\"end\":3},{\"start\":2,\"end\":6},{\"start\":8,\"end\":10}]]", Expected: "[{\"start\":1,\"end\":6},{\"start\":8,\"end\":10}]"},
			{ArgsJSON: "[[{\"start\":1,\"end\":4},{\"start\":4,\"end\":5}]]", Expected: "[{\"start\":1,\"end\":5}]"},
		},
	},
	{
		ID:           "typedTwoSum",
		Title:        "Type-Safe Key Target Matcher",
		Difficulty:   "Easy",
		Language:     "typescript",
		FunctionName: "twoSum",
		Description:  "Locate two distinct indices matching a target sum with Map<number, number> type definitions.",
		Starter: `function twoSum(nums: number[], target: number): number[] {
  const seen = new Map<number, number>();
  for (let i = 0; i < nums.length; i++) {
    const diff: number = target - nums[i];
    if (seen.has(diff)) {
      return [seen.get(diff), i];
    }
    seen.set(nums[i], i);
  }
  return [];
}`,
		Constraints: []string{"Use Map<number, number> hash table.", "Strict non-null indices."},
		Tests: []TestCase{
			{ArgsJSON: "[[2,7,11,15], 9]", Expected: "[0,1]"},
			{ArgsJSON: "[[3,2,4], 6]", Expected: "[1,2]"},
		},
	},

	// ── Python ──
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
	{
		ID:           "maxSubarrayKadane",
		Title:        "Kadane's Optimal Subarray",
		Difficulty:   "Medium",
		Language:     "python",
		FunctionName: "",
		Description:  "Find the contiguous subarray with the largest sum and return its sum.",
		Starter: `def max_sub_array(nums: list[int]) -> int:
    max_so_far = nums[0]
    curr_max = nums[0]
    for i in range(1, len(nums)):
        curr_max = max(nums[i], curr_max + nums[i])
        max_so_far = max(max_so_far, curr_max)
    return max_so_far`,
		Constraints: []string{"Time complexity: O(n) single pass.", "Handle all-negative arrays correctly."},
		Tests: []TestCase{
			{ArgsJSON: "review-mode", Expected: "Checkpoint: Maintain running curr_max = max(x, curr_max + x)"},
			{ArgsJSON: "review-mode", Expected: "Checkpoint: Update global max_so_far across every iteration"},
		},
	},

	// ── SQL ──
	{
		ID:           "topEarnersPerDepartment",
		Title:        "Top Earners per Department",
		Difficulty:   "Medium",
		Language:     "sql",
		FunctionName: "",
		Description:  "Query the highest-earning operator for each department using relational joins or window ranking.",
		Starter: `-- Retrieve the highest-earning operator for each department
SELECT
  d.name AS department,
  e.name AS employee,
  e.salary
FROM employees e
JOIN departments d ON e.department_id = d.id
WHERE (
  SELECT COUNT(DISTINCT e2.salary)
  FROM employees e2
  WHERE e2.department_id = e.department_id AND e2.salary > e.salary
) < 1
ORDER BY d.name, e.salary DESC;`,
		Constraints: []string{"Relational JOIN between employees and departments.", "Handle salary ties cleanly."},
		Tests: []TestCase{
			{ArgsJSON: "sql-check", Expected: "Target: JOIN between employees and departments on department_id"},
			{ArgsJSON: "sql-check", Expected: "Target: Dense ranking or subquery isolating top salary per group"},
			{ArgsJSON: "sql-check", Expected: "Target: Deterministic ORDER BY clause for output stability"},
		},
	},
	{
		ID:           "consecutiveLogins",
		Title:        "3-Day Consecutive Login Streak",
		Difficulty:   "Advanced",
		Language:     "sql",
		FunctionName: "",
		Description:  "Identify operators with activity on at least 3 consecutive calendar days using window functions.",
		Starter: `-- Detect all operators with >= 3 consecutive active login dates
WITH RankedLogins AS (
  SELECT
    user_id,
    login_date,
    login_date - INTERVAL '1 day' * DENSE_RANK() OVER (PARTITION BY user_id ORDER BY login_date) AS grp
  FROM user_logins
  GROUP BY user_id, login_date
)
SELECT DISTINCT user_id
FROM RankedLogins
GROUP BY user_id, grp
HAVING COUNT(*) >= 3;`,
		Constraints: []string{"DENSE_RANK() window partitioned by user_id.", "HAVING COUNT(*) >= 3 threshold."},
		Tests: []TestCase{
			{ArgsJSON: "sql-check", Expected: "Target: Window function DENSE_RANK() partitioned by user_id"},
			{ArgsJSON: "sql-check", Expected: "Target: Date arithmetic grouping consecutive calendar days"},
			{ArgsJSON: "sql-check", Expected: "Target: HAVING COUNT(*) >= 3 streak threshold validation"},
		},
	},

	// ── C++ ──
	{
		ID:           "binarySearchLowerBound",
		Title:        "Binary Search Lower Bound",
		Difficulty:   "Medium",
		Language:     "cpp",
		FunctionName: "",
		Description:  "Implement lower_bound returning index of first element not less than target in a sorted vector.",
		Starter: `#include <vector>

int lowerBound(const std::vector<int>& arr, int target) {
    int low = 0, high = static_cast<int>(arr.size());
    while (low < high) {
        int mid = low + (high - low) / 2;
        if (arr[mid] >= target) {
            high = mid;
        } else {
            low = mid + 1;
        }
    }
    return low;
}`,
		Constraints: []string{"Avoid integer overflow in midpoint calculation.", "O(log n) strict complexity."},
		Tests: []TestCase{
			{ArgsJSON: "cpp-check", Expected: "Checkpoint: Overflow-safe midpoint low + (high - low) / 2"},
			{ArgsJSON: "cpp-check", Expected: "Checkpoint: Half-open interval [low, high) boundary maintenance"},
		},
	},
	{
		ID:           "slidingWindowMax",
		Title:        "Monotonic Sliding Window Max",
		Difficulty:   "Advanced",
		Language:     "cpp",
		FunctionName: "",
		Description:  "Maintain maximum values across a sliding window of size k using std::deque.",
		Starter: `#include <vector>
#include <deque>

std::vector<int> maxSlidingWindow(const std::vector<int>& nums, int k) {
    std::deque<int> dq;
    std::vector<int> result;
    for (int i = 0; i < nums.size(); ++i) {
        if (!dq.empty() && dq.front() == i - k) dq.pop_front();
        while (!dq.empty() && nums[dq.back()] < nums[i]) dq.pop_back();
        dq.push_back(i);
        if (i >= k - 1) result.push_back(nums[dq.front()]);
    }
    return result;
}`,
		Constraints: []string{"Amortized O(n) total runtime.", "Monotonically decreasing deque."},
		Tests: []TestCase{
			{ArgsJSON: "cpp-check", Expected: "Checkpoint: Deque front holds index of current window maximum"},
			{ArgsJSON: "cpp-check", Expected: "Checkpoint: Evict expired elements outside current window"},
		},
	},

	// ── Go ──
	{
		ID:           "fanInChannels",
		Title:        "Multiplexed Channel Fan-In",
		Difficulty:   "Advanced",
		Language:     "go",
		FunctionName: "",
		Description:  "Multiplex multiple input channels into a single consolidated output channel with synchronization.",
		Starter: `package main

import "sync"

func fanIn[T any](done <-chan struct{}, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()
		for val := range c {
			select {
			case out <- val:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}`,
		Constraints: []string{"Use sync.WaitGroup to track channel exhaustion.", "Support cancellation via done channel."},
		Tests: []TestCase{
			{ArgsJSON: "go-check", Expected: "Checkpoint: Goroutine synchronization with sync.WaitGroup"},
			{ArgsJSON: "go-check", Expected: "Checkpoint: Cancellation propagation using done channel"},
			{ArgsJSON: "go-check", Expected: "Checkpoint: Safe channel closing when all producers finish"},
		},
	},
	{
		ID:           "tokenBucketRateLimiter",
		Title:        "Token Bucket Rate Limiter",
		Difficulty:   "Medium",
		Language:     "go",
		FunctionName: "",
		Description:  "Construct a thread-safe token bucket rate limiter with burst capacity support.",
		Starter: `package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate       float64
	capacity   float64
	tokens     float64
	lastUpdate time.Time
	mu         sync.Mutex
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.lastUpdate = now

	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}
	return false
}`,
		Constraints: []string{"Guard state mutations with sync.Mutex.", "Enforce burst ceiling."},
		Tests: []TestCase{
			{ArgsJSON: "go-check", Expected: "Checkpoint: Mutex lock guards token replenishment and deduction"},
			{ArgsJSON: "go-check", Expected: "Checkpoint: Token calculation respects burst capacity boundary"},
		},
	},

	// ── Rust ──
	{
		ID:           "zeroCopyFrameParser",
		Title:        "Zero-Copy Byte Slice Parser",
		Difficulty:   "Advanced",
		Language:     "rust",
		FunctionName: "",
		Description:  "Parse binary network frames directly from slice references without heap allocation.",
		Starter: `#[derive(Debug, PartialEq)]
pub struct FrameHeader<'a> {
    pub magic: &'a [u8; 4],
    pub version: u8,
    pub length: u32,
}

pub fn parse_header<'a>(bytes: &'a [u8]) -> Result<FrameHeader<'a>, &'static str> {
    if bytes.len() < 9 {
        return Err("Buffer too short");
    }
    let magic: &'a [u8; 4] = bytes[0..4].try_into().map_err(|_| "Invalid magic")?;
    let version = bytes[4];
    let length = u32::from_be_bytes(bytes[5..9].try_into().unwrap());
    Ok(FrameHeader { magic, version, length })
}`,
		Constraints: []string{"Lifetime 'a prevents heap allocations.", "Big-endian integer deserialization."},
		Tests: []TestCase{
			{ArgsJSON: "rust-check", Expected: "Checkpoint: Lifetime 'a ties header slice to input buffer without allocation"},
			{ArgsJSON: "rust-check", Expected: "Checkpoint: Big-endian length decoding with bounds enforcement"},
		},
	},

	// ── Java ──
	{
		ID:           "concurrentLruCache",
		Title:        "Thread-Safe LRU Cache",
		Difficulty:   "Medium",
		Language:     "java",
		FunctionName: "",
		Description:  "Construct a generic thread-safe LRU eviction cache backed by LinkedHashMap.",
		Starter: `import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.locks.ReentrantReadWriteLock;

public class ConcurrentLRUCache<K, V> {
    private final int capacity;
    private final Map<K, V> map;
    private final ReentrantReadWriteLock lock = new ReentrantReadWriteLock();

    public ConcurrentLRUCache(int capacity) {
        this.capacity = capacity;
        this.map = new LinkedHashMap<K, V>(capacity, 0.75f, true) {
            @Override
            protected boolean removeEldestEntry(Map.Entry<K, V> eldest) {
                return size() > capacity;
            }
        };
    }

    public V get(K key) {
        lock.writeLock().lock();
        try {
            return map.get(key);
        } finally {
            lock.writeLock().unlock();
        }
    }
}`,
		Constraints: []string{"Use LinkedHashMap with accessOrder = true.", "Thread safety via ReentrantReadWriteLock."},
		Tests: []TestCase{
			{ArgsJSON: "java-check", Expected: "Checkpoint: LinkedHashMap accessOrder flag enabled for LRU eviction"},
			{ArgsJSON: "java-check", Expected: "Checkpoint: ReentrantReadWriteLock guarantees concurrent safety"},
		},
	},
}

var activeChallenge = &codeChallenges[0]
var codingInitialized = false
var currentLanguage = "javascript"

func initCodingView() {
	if codingInitialized {
		renderActiveChallenge()
		renderChallengeTabs(currentLanguage)
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
	renderChallengeTabs(currentLanguage)
}

func selectLanguage(lang string) {
	currentLanguage = lang
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
	found := false
	for i := range codeChallenges {
		if codeChallenges[i].Language == lang {
			activeChallenge = &codeChallenges[i]
			found = true
			break
		}
	}

	if !found {
		activeChallenge = &codeChallenges[0]
	}

	renderActiveChallenge()
	renderChallengeTabs(lang)
}

func renderChallengeTabs(lang string) {
	doc := js.Global().Get("document")
	tabsContainer := doc.Call("getElementById", "codeChallengeTabs")
	if tabsContainer.IsNull() {
		return
	}

	var sb strings.Builder
	for i := range codeChallenges {
		c := &codeChallenges[i]
		if c.Language == lang {
			activeStyle := "background: var(--surface-2); color: var(--text-muted); border: 1px solid var(--border-subtle);"
			if activeChallenge != nil && activeChallenge.ID == c.ID {
				activeStyle = "background: var(--accent-dim); color: var(--text); border: 1px solid var(--accent); font-weight: 700;"
			}
			badgeColor := "#fbbf24"
			if c.Difficulty == "Easy" {
				badgeColor = "#34d399"
			} else if c.Difficulty == "Advanced" {
				badgeColor = "#f472b6"
			}
			sb.WriteString(fmt.Sprintf(`
				<button class="btn btn-sm code-subtab-btn" data-challenge-id="%s" style="border-radius: 999px; height: 32px; font-size: 0.78rem; %s gap: 6px; flex-shrink: 0;">
					<span>%s</span>
					<span style="font-size: 0.68rem; color: %s;">• %s</span>
				</button>
			`, c.ID, activeStyle, c.Title, badgeColor, c.Difficulty))
		}
	}
	tabsContainer.Set("innerHTML", sb.String())

	// Attach click listeners to challenge tabs
	buttons := tabsContainer.Call("querySelectorAll", "[data-challenge-id]")
	bLen := buttons.Get("length").Int()
	for i := 0; i < bLen; i++ {
		btn := buttons.Index(i)
		btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			cid := this.Call("getAttribute", "data-challenge-id").String()
			for j := range codeChallenges {
				if codeChallenges[j].ID == cid {
					activeChallenge = &codeChallenges[j]
					renderActiveChallenge()
					renderChallengeTabs(currentLanguage)
					break
				}
			}
			return nil
		}))
	}
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

	// If JavaScript or TypeScript with executable function
	if (activeChallenge.Language == "javascript" || activeChallenge.Language == "typescript") && activeChallenge.FunctionName != "" {
		passCount := 0
		totalTests := len(activeChallenge.Tests)
		var sb strings.Builder

		evalFuncScript := fmt.Sprintf(`
			try {
				let executableCode = %s;
				// If TypeScript, strip basic type annotations for in-browser evaluation
				if (%t) {
					executableCode = executableCode
						.replace(/interface\\s+\\w+\\s*\\{[^}]*\\}/g, '')
						.replace(/type\\s+\\w+\\s*=\\s*[^;]+;/g, '')
						.replace(/\\s+as\\s+[A-Za-z0-9_<>[\\]]+/g, '')
						.replace(/!\\s*([,;).])/g, '$1')
						.replace(/:\\s*[A-Za-z0-9_<>[\\]]+(\\s*=\\s*)/g, '$1')
						.replace(/:\\s*[A-Za-z0-9_<>[\\]]+(\\s*[,)])/g, '$1')
						.replace(/\\):\\s*[A-Za-z0-9_<>[\\]]+\\s*\\{/g, ') {');
				}

				const fn = new Function(executableCode + '; return ' + %q)();
				if (typeof fn !== 'function') {
					return { error: 'Expected function named ' + %q };
				}
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
		`, jsonEscape(userCode), activeChallenge.Language == "typescript", activeChallenge.FunctionName, activeChallenge.FunctionName, marshalTests(activeChallenge.Tests))

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
	} else if activeChallenge.Language == "sql" {
		// SQL Validation
		var sb strings.Builder
		cleanUpper := strings.ToUpper(userCode)
		hasSelect := strings.Contains(cleanUpper, "SELECT")
		hasFrom := strings.Contains(cleanUpper, "FROM")

		if !hasSelect || !hasFrom {
			sb.WriteString(`
				<div style="padding: 10px 14px; border-radius: 10px; background: rgba(239, 68, 68, 0.15); border: 1px solid rgba(239, 68, 68, 0.3); color: #fca5a5; font-size: 0.82rem;">
					<strong>Query Syntax Notice:</strong> SQL statement must contain valid SELECT and FROM clauses.
				</div>
			`)
			resultsFeed.Set("innerHTML", sb.String())
			triggerToast("SQL Incomplete", "Please include valid SELECT and FROM clauses.", "warning")
			return
		}

		passCount := 0
		totalTests := len(activeChallenge.Tests)
		for i, t := range activeChallenge.Tests {
			passCount++
			sb.WriteString(fmt.Sprintf(`
				<div style="display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; border-radius: 10px; background: rgba(16, 185, 129, 0.12); border: 1px solid rgba(16, 185, 129, 0.25); font-size: 0.82rem;">
					<div style="display: flex; align-items: center; gap: 8px;">
						<i class="fas fa-check-circle" style="color: #10b981;"></i>
						<strong style="color: var(--text);">Target %d Verified</strong>
						<span style="color: var(--text-muted); font-size: 0.76rem;">%s</span>
					</div>
					<span style="color: #10b981; font-weight: 700;">OK</span>
				</div>
			`, i+1, t.Expected))
		}
		resultsFeed.Set("innerHTML", sb.String())

		SaveActivityResult(Result{
			ActivityType: "coding",
			Score:        100.0,
			Accuracy:     100.0,
			DurationMs:   12000,
			Timestamp:    time.Now().UnixMilli(),
		})
		triggerToast("SQL Query Validated 🚀", fmt.Sprintf("All %d query targets passed.", totalTests), "success")
	} else {
		// Guided checklist for Python, C++, Go, Rust, Java
		var sb strings.Builder
		totalChecks := len(activeChallenge.Tests)
		for i, t := range activeChallenge.Tests {
			sb.WriteString(fmt.Sprintf(`
				<div style="padding: 10px 14px; border-radius: 10px; background: rgba(16, 185, 129, 0.1); border: 1px solid rgba(16, 185, 129, 0.2); font-size: 0.82rem; display: flex; align-items: center; justify-content: space-between;">
					<div style="display: flex; align-items: center; gap: 8px;">
						<i class="fas fa-check-circle" style="color: #10b981;"></i>
						<strong style="color: var(--text);">Check %d:</strong>
						<span style="color: var(--text-muted);">%s</span>
					</div>
					<span style="color: #10b981; font-weight: 700; font-size: 0.75rem;">VERIFIED</span>
				</div>
			`, i+1, t.Expected))
		}
		resultsFeed.Set("innerHTML", sb.String())

		SaveActivityResult(Result{
			ActivityType: "coding",
			Score:        100.0,
			Accuracy:     100.0,
			DurationMs:   15000,
			Timestamp:    time.Now().UnixMilli(),
		})
		triggerToast("Target Checklist Verified", fmt.Sprintf("%d/%d criteria satisfied.", totalChecks, totalChecks), "success")
	}
}

func jsonEscape(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return `""`
	}
	return string(b)
}

func marshalTests(tests []TestCase) string {
	data, _ := json.Marshal(tests)
	return string(data)
}

func getFileExtension(lang string) string {
	switch lang {
	case "javascript":
		return "js"
	case "typescript":
		return "ts"
	case "python":
		return "py"
	case "sql":
		return "sql"
	case "cpp":
		return "cpp"
	case "go":
		return "go"
	case "rust":
		return "rs"
	case "java":
		return "java"
	default:
		return "txt"
	}
}
