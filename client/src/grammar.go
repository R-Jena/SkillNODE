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

type GrammarQuestion struct {
	ID          string
	Category    string
	Prompt      string
	Options     []string
	AnswerIndex int
	Explanation string
}

var grammarBank = []GrammarQuestion{
	{
		ID:          "subj-verb-neither",
		Category:    "Subject-Verb Agreement",
		Prompt:      "Neither the distributed storage nodes nor the primary coordinator ________ responding to health probes.",
		Options:     []string{"is", "are", "were", "being"},
		AnswerIndex: 0,
		Explanation: "When subject elements are joined by 'neither... nor', the verb agrees in number with the nearer subject ('primary coordinator' is singular -> 'is').",
	},
	{
		ID:          "parallelism-verbs",
		Category:    "Parallel Structure",
		Prompt:      "The ingest pipeline is configured for validating payloads, archiving snapshots, and ________ telemetry alerts.",
		Options:     []string{"broadcasting", "to broadcast", "will broadcast", "broadcasted"},
		AnswerIndex: 0,
		Explanation: "Items in a coordinate series must share identical grammatical form (gerunds: validating, archiving, broadcasting).",
	},
	{
		ID:          "dangling-modifier-compile",
		Category:    "Modifiers",
		Prompt:      "Choose the sentence free of dangling or misplaced modifiers:",
		Options:     []string{
			"Having compiled the Go WebAssembly core, the browser executed the binary immediately.",
			"Having compiled the Go WebAssembly core, the developer verified immediate execution.",
			"Compiling the Go WebAssembly core, the binary ran without memory overhead.",
			"When compiled, zero latency was observed by the benchmark engine.",
		},
		AnswerIndex: 1,
		Explanation: "The introductory participial phrase 'Having compiled...' must logically modify the agent who compiled it ('the developer').",
	},
	{
		ID:          "subjunctive-mandative",
		Category:    "Subjunctive Mood",
		Prompt:      "The architecture committee recommended that the authentication token ________ rotated every thirty days.",
		Options:     []string{"be", "is", "was", "must be"},
		AnswerIndex: 0,
		Explanation: "Mandative clauses following verbs of recommendation or demand ('recommended that...') require the subjunctive base form ('be').",
	},
	{
		ID:          "fewer-vs-less",
		Category:    "Diction & Precision",
		Prompt:      "The optimized hash index requires ________ memory allocations per query than the previous tree implementation.",
		Options:     []string{"fewer", "less", "lesser", "fewest"},
		AnswerIndex: 0,
		Explanation: "Use 'fewer' for countable discrete entities ('allocations') and 'less' for continuous uncountable quantities ('memory space').",
	},
	{
		ID:          "who-vs-whom",
		Category:    "Pronoun Case",
		Prompt:      "The systems architect ________ the board appointed to supervise the migration delivered zero downtime.",
		Options:     []string{"whom", "who", "which", "whose"},
		AnswerIndex: 0,
		Explanation: "Use 'whom' as the objective pronoun when it serves as the grammatical object of the clause verb ('the board appointed [whom]').",
	},
	{
		ID:          "affect-vs-effect",
		Category:    "Diction & Precision",
		Prompt:      "Elevated latency spikes in downstream clusters directly ________ overall transactional throughput.",
		Options:     []string{"affect", "effect", "affected", "effecting"},
		AnswerIndex: 0,
		Explanation: "'Affect' is typically a verb meaning to produce an influence upon; 'effect' is typically a noun denoting the result.",
	},
	{
		ID:          "pronoun-antecedent-every",
		Category:    "Pronoun-Antecedent Agreement",
		Prompt:      "Every service instance must flush ________ local telemetry queue before shutdown.",
		Options:     []string{"its", "their", "it's", "there"},
		AnswerIndex: 0,
		Explanation: "'Every service instance' is grammatically singular, requiring the singular possessive pronoun 'its' (without apostrophe).",
	},
	{
		ID:          "semicolon-conjunctive-adverb",
		Category:    "Punctuation & Mechanics",
		Prompt:      "The primary database node experienced kernel panic; ________, the replica promoted seamlessly without dropped transactions.",
		Options:     []string{"however", "therefore", "because", "although"},
		AnswerIndex: 0,
		Explanation: "Use a semicolon followed by a conjunctive adverb like 'however' (with trailing comma) to contrast two independent coordinate clauses.",
	},
	{
		ID:          "redundancy-elimination",
		Category:    "Conciseness & Style",
		Prompt:      "Identify the sentence free of redundant phrasing:",
		Options:     []string{
			"The engineering teams must collaborate together on the schema revision.",
			"Advance planning is mandatory prior to upgrading database clusters.",
			"The teams must collaborate on the schema revision.",
			"In my personal opinion, we arrived at a consensus of opinion.",
		},
		AnswerIndex: 2,
		Explanation: "'Collaborate' already implies togetherness; 'advance planning' and 'consensus of opinion' are pleonastic redundancies.",
	},
	{
		ID:          "discrete-vs-discreet",
		Category:    "Diction & Precision",
		Prompt:      "The microservices communicate across twelve ________ network namespaces.",
		Options:     []string{"discrete", "discreet", "discreted", "discretionary"},
		AnswerIndex: 0,
		Explanation: "'Discrete' means distinct, detached, or separate; 'discreet' means cautious, prudent, or unobtrusive.",
	},
	{
		ID:          "comma-splice-identification",
		Category:    "Punctuation & Mechanics",
		Prompt:      "Select the grammatically valid sentence:",
		Options:     []string{
			"The API Gateway rejected the payload, the HMAC signature was invalid.",
			"The API Gateway rejected the payload; the HMAC signature was invalid.",
			"The API Gateway rejected the payload, because the signature was invalid.",
			"The API Gateway rejected the payload, but the HMAC signature was invalid.",
		},
		AnswerIndex: 1,
		Explanation: "Two independent clauses cannot be joined solely by a comma without a coordinating conjunction (comma splice). A semicolon correctly links them.",
	},
	{
		ID:          "third-conditional",
		Category:    "Conditionals & Tenses",
		Prompt:      "If the cluster coordinator ________ the heartbeat signal, failover would not have occurred.",
		Options:     []string{"had received", "would have received", "received", "has received"},
		AnswerIndex: 0,
		Explanation: "Past counterfactual conditional clauses (Third Conditional) employ past perfect ('had received') in the if-clause.",
	},
	{
		ID:          "principal-vs-principle",
		Category:    "Diction & Precision",
		Prompt:      "Idempotency is a foundational ________ of robust distributed system architecture.",
		Options:     []string{"principle", "principal", "principled", "principality"},
		AnswerIndex: 0,
		Explanation: "A 'principle' is a fundamental truth or law; 'principal' refers to a chief person, main sum, or primary element.",
	},
	{
		ID:          "collective-noun-concord",
		Category:    "Subject-Verb Agreement",
		Prompt:      "The site reliability committee ________ unanimously approved the production rollout.",
		Options:     []string{"has", "have", "are", "were"},
		AnswerIndex: 0,
		Explanation: "Collective nouns acting as a unified singular entity ('the committee') take singular verbs ('has approved').",
	},
	{
		ID:          "compliment-vs-complement",
		Category:    "Diction & Precision",
		Prompt:      "The asynchronous worker pool serves as an ideal ________ to the synchronous request-response pipeline.",
		Options:     []string{"complement", "compliment", "complementary", "complimentary"},
		AnswerIndex: 0,
		Explanation: "'Complement' means something that completes or enhances; 'compliment' denotes praise or flattery.",
	},
	{
		ID:          "correlative-parallelism",
		Category:    "Parallel Structure",
		Prompt:      "The service guarantees not only low latency ________ high durability across partitions.",
		Options:     []string{"but also", "and also", "as well as", "along with"},
		AnswerIndex: 0,
		Explanation: "The correlative conjunction pair is 'not only... but also', requiring parallel structural alignment.",
	},
	{
		ID:          "criterion-vs-criteria",
		Category:    "Pluralization & Concord",
		Prompt:      "Latency under 10 milliseconds was the sole ________ for production certification.",
		Options:     []string{"criterion", "criteria", "criterions", "criterium"},
		AnswerIndex: 0,
		Explanation: "'Criterion' is the singular noun; 'criteria' is plural. Since the sentence specifies 'the sole', singular 'criterion' is required.",
	},
	{
		ID:          "preposition-accordance",
		Category:    "Idiomatic Prepositions",
		Prompt:      "The cryptographic envelope was generated in accordance ________ zero-trust specifications.",
		Options:     []string{"with", "to", "by", "of"},
		AnswerIndex: 0,
		Explanation: "The correct idiomatic English construction is 'in accordance with'.",
	},
	{
		ID:          "active-voice-preference",
		Category:    "Conciseness & Style",
		Prompt:      "Choose the most concise active-voice revision for: 'A telemetry packet was dispatched by the client.'",
		Options:     []string{
			"The client dispatched a telemetry packet.",
			"Dispatched by the client was a telemetry packet.",
			"A telemetry packet was being dispatched by the client.",
			"By the client, a telemetry packet was dispatched.",
		},
		AnswerIndex: 0,
		Explanation: "'The client dispatched a telemetry packet' eliminates passive auxiliary verbs and places the grammatical agent upfront.",
	},
	{
		ID:          "irregardless-nonstandard",
		Category:    "Standard Usage",
		Prompt:      "The deployment continued ________ of downstream latency warnings.",
		Options:     []string{"regardless", "irregardless", "disregardless", "unregarding"},
		AnswerIndex: 0,
		Explanation: "'Regardless' is standard English; 'irregardless' is a nonstandard double negative tautology.",
	},
	{
		ID:          "lie-vs-lay",
		Category:    "Diction & Precision",
		Prompt:      "The backup archives ________ uncompressed on the cold storage bucket.",
		Options:     []string{"lie", "lay", "laid", "lain"},
		AnswerIndex: 0,
		Explanation: "'Lie' (present tense, intransitive) means to recline or rest; 'lay' (transitive) requires a direct object ('lay the archives down').",
	},
	{
		ID:          "em-dash-mechanics",
		Category:    "Punctuation & Mechanics",
		Prompt:      "Select the sentence using em-dashes correctly for parenthetical emphasis:",
		Options:     []string{
			"The core runtime—written entirely in Go WebAssembly—executes with sub-millisecond overhead.",
			"The core runtime—written entirely in Go WebAssembly, executes with sub-millisecond overhead.",
			"The core runtime, written entirely in Go WebAssembly—executes with sub-millisecond overhead.",
			"The core runtime - written entirely in Go WebAssembly - executes with sub-millisecond overhead.",
		},
		AnswerIndex: 0,
		Explanation: "Paired em-dashes must enclose parenthetical insertions symmetrically without trailing commas.",
	},
	{
		ID:          "each-verb-agreement",
		Category:    "Subject-Verb Agreement",
		Prompt:      "Each of the twenty worker threads ________ assigned an independent execution context.",
		Options:     []string{"is", "are", "were", "have been"},
		AnswerIndex: 0,
		Explanation: "The pronoun 'Each' is grammatically singular and governs the verb, regardless of plural prepositional complements ('of the twenty worker threads').",
	},
	{
		ID:          "as-if-subjunctive",
		Category:    "Subjunctive Mood",
		Prompt:      "The load balancer behaves as though all downstream sockets ________ saturated.",
		Options:     []string{"were", "was", "are", "is"},
		AnswerIndex: 0,
		Explanation: "Clauses introduced by 'as if' or 'as though' describing hypothetical or counterfactual states take subjunctive 'were'.",
	},
	{
		ID:          "its-vs-its-apostrophe",
		Category:    "Punctuation & Mechanics",
		Prompt:      "The runtime checked ________ memory bounds before executing the instruction.",
		Options:     []string{"its", "it's", "its'", "its's"},
		AnswerIndex: 0,
		Explanation: "'Its' is the possessive pronoun; 'it's' is a contraction for 'it is' or 'it has'.",
	},
	{
		ID:          "phenomenon-plural",
		Category:    "Pluralization & Concord",
		Prompt:      "Network partitioning is an unavoidable ________ in wide-area distributed computing.",
		Options:     []string{"phenomenon", "phenomena", "phenomenons", "phenomenal"},
		AnswerIndex: 0,
		Explanation: "'Phenomenon' is singular; 'phenomena' is plural. 'An unavoidable [singular]' demands 'phenomenon'.",
	},
	{
		ID:          "than-case",
		Category:    "Pronoun Case",
		Prompt:      "The senior kernel architect benchmarked more iterations than ________.",
		Options:     []string{"I", "me", "myself", "mine"},
		AnswerIndex: 0,
		Explanation: "In elliptical comparative constructions ('than I [benchmarked]'), the subjective case pronoun 'I' is formally correct.",
	},
	{
		ID:          "between-vs-among",
		Category:    "Diction & Precision",
		Prompt:      "Workload distribution was partitioned evenly ________ the five database replicas.",
		Options:     []string{"among", "between", "amid", "along"},
		AnswerIndex: 0,
		Explanation: "Use 'among' when referring to collective relationships among three or more entities without distinct one-to-one pairs.",
	},
	{
		ID:          "insure-vs-ensure",
		Category:    "Diction & Precision",
		Prompt:      "Strict checksum validation was implemented to ________ data integrity across network hops.",
		Options:     []string{"ensure", "insure", "assure", "reassure"},
		AnswerIndex: 0,
		Explanation: "'Ensure' means to make certain; 'insure' refers to financial insurance underwriting; 'assure' means to give confidence to someone.",
	},
}

type GrammarState struct {
	Running       bool
	CurrentIndex  int
	Score         int
	Streak        int
	Multiplier    float64
	AnsweredCount int
	CorrectCount  int
	SelectedOrder []int
}

var gramState = &GrammarState{
	Multiplier: 1.0,
}
var grammarInitialized = false

func initGrammarView() {
	if grammarInitialized {
		return
	}
	grammarInitialized = true

	doc := js.Global().Get("document")

	startBtn := doc.Call("getElementById", "grammarStartBtn")
	if !startBtn.IsNull() {
		startBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			startGrammarRound()
			return nil
		}))
	}

	nextBtn := doc.Call("getElementById", "grammarNextBtn")
	if !nextBtn.IsNull() {
		nextBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			nextGrammarQuestion()
			return nil
		}))
	}

	resetGrammarState()
}

func resetGrammarState() {
	gramState.Running = false
	gramState.Score = 0
	gramState.Streak = 0
	gramState.Multiplier = 1.0
	gramState.AnsweredCount = 0
	gramState.CorrectCount = 0
	gramState.CurrentIndex = 0

	// Shuffle questions order
	rand.Seed(time.Now().UnixNano())
	n := len(grammarBank)
	gramState.SelectedOrder = rand.Perm(n)

	updateGrammarHUD()

	doc := js.Global().Get("document")
	promptEl := doc.Call("getElementById", "grammarPrompt")
	if !promptEl.IsNull() {
		promptEl.Set("innerText", "Click 'Start Session' to test grammatical discernment and syntactical precision.")
	}

	optionsEl := doc.Call("getElementById", "grammarOptionsContainer")
	if !optionsEl.IsNull() {
		optionsEl.Set("innerHTML", "")
	}

	feedbackEl := doc.Call("getElementById", "grammarFeedbackNote")
	if !feedbackEl.IsNull() {
		feedbackEl.Set("innerHTML", `<div style="font-size: 0.85rem; color: var(--text-faint);">30+ curated questions covering conciseness, agreement, subjunctive clauses, and technical mechanics.</div>`)
	}

	btnNext := doc.Call("getElementById", "grammarNextBtn")
	if !btnNext.IsNull() {
		btnNext.Get("style").Set("display", "none")
	}
}

func startGrammarRound() {
	resetGrammarState()
	gramState.Running = true
	renderGrammarQuestion()
}

func renderGrammarQuestion() {
	if gramState.CurrentIndex >= len(gramState.SelectedOrder) {
		finishGrammarRound()
		return
	}

	qIdx := gramState.SelectedOrder[gramState.CurrentIndex]
	q := grammarBank[qIdx]

	doc := js.Global().Get("document")

	setElementText(doc, "grammarCategoryBadge", q.Category)
	setElementText(doc, "grammarPrompt", q.Prompt)
	setElementText(doc, "grammarProgressText", fmt.Sprintf("Question %d of %d", gramState.CurrentIndex+1, len(gramState.SelectedOrder)))

	btnNext := doc.Call("getElementById", "grammarNextBtn")
	if !btnNext.IsNull() {
		btnNext.Get("style").Set("display", "none")
	}

	feedbackEl := doc.Call("getElementById", "grammarFeedbackNote")
	if !feedbackEl.IsNull() {
		feedbackEl.Set("innerHTML", `<div style="font-size: 0.82rem; color: var(--text-faint);">Select the grammatically sound replacement.</div>`)
	}

	// Render options
	optionsContainer := doc.Call("getElementById", "grammarOptionsContainer")
	if !optionsContainer.IsNull() {
		var sb strings.Builder
		sb.WriteString(`<div style="display: flex; flex-direction: column; gap: 10px; width: 100%;">`)
		for i, opt := range q.Options {
			sb.WriteString(fmt.Sprintf(`
				<button class="btn btn-secondary grammar-opt-btn" data-opt-idx="%d" style="justify-content: flex-start; text-align: left; padding: 14px 18px; height: auto; border-radius: 12px; font-size: 0.92rem; line-height: 1.5; white-space: normal;">
					<span style="display: inline-flex; width: 24px; height: 24px; align-items: center; justify-content: center; border-radius: 6px; background: var(--surface-3); font-weight: 700; margin-right: 12px; font-size: 0.78rem;">%c</span>
					<span>%s</span>
				</button>
			`, i, 'A'+i, opt))
		}
		sb.WriteString(`</div>`)
		optionsContainer.Set("innerHTML", sb.String())

		// Bind option clicks
		btns := optionsContainer.Call("querySelectorAll", ".grammar-opt-btn")
		bCount := btns.Get("length").Int()
		for i := 0; i < bCount; i++ {
			btn := btns.Index(i)
			btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				chosen := this.Call("getAttribute", "data-opt-idx").Int()
				evaluateGrammarAnswer(chosen)
				return nil
			}))
		}
	}
}

func evaluateGrammarAnswer(chosenIdx int) {
	if !gramState.Running || gramState.CurrentIndex >= len(gramState.SelectedOrder) {
		return
	}

	qIdx := gramState.SelectedOrder[gramState.CurrentIndex]
	q := grammarBank[qIdx]
	doc := js.Global().Get("document")

	// Disable buttons
	btns := doc.Call("querySelectorAll", ".grammar-opt-btn")
	bCount := btns.Get("length").Int()
	for i := 0; i < bCount; i++ {
		btns.Index(i).Set("disabled", true)
	}

	gramState.AnsweredCount++
	isCorrect := (chosenIdx == q.AnswerIndex)

	pts := 0
	if isCorrect {
		gramState.CorrectCount++
		gramState.Streak++
		if gramState.Streak >= 10 {
			gramState.Multiplier = 3.0
		} else if gramState.Streak >= 6 {
			gramState.Multiplier = 2.0
		} else if gramState.Streak >= 3 {
			gramState.Multiplier = 1.5
		}
		pts = int(100.0 * gramState.Multiplier)
		gramState.Score += pts
		triggerToast("Correct! 🎯", fmt.Sprintf("+%d Points (Streak: %d 🔥)", pts, gramState.Streak), "success")
	} else {
		gramState.Streak = 0
		gramState.Multiplier = 1.0
		triggerToast("Incorrect", "Review the grammatical rationale below.", "warning")
	}

	updateGrammarHUD()

	// Show feedback and explanation
	feedbackEl := doc.Call("getElementById", "grammarFeedbackNote")
	if !feedbackEl.IsNull() {
		statusColor := "#10b981"
		statusTitle := "Correct Choice"
		if !isCorrect {
			statusColor = "#ef4444"
			statusTitle = "Incorrect"
		}

		html := fmt.Sprintf(`
			<div style="padding: 16px; border-radius: 12px; background: var(--surface-2); border: 1px solid var(--border-subtle); margin-top: 12px;">
				<div style="font-weight: 700; color: %s; font-size: 0.95rem; margin-bottom: 6px;">
					%s
				</div>
				<div style="font-size: 0.85rem; color: var(--text); line-height: 1.6;">
					<strong>Grammar Rule:</strong> %s
				</div>
			</div>
		`, statusColor, statusTitle, q.Explanation)
		feedbackEl.Set("innerHTML", html)
	}

	btnNext := doc.Call("getElementById", "grammarNextBtn")
	if !btnNext.IsNull() {
		btnNext.Get("style").Set("display", "inline-flex")
	}
}

func nextGrammarQuestion() {
	gramState.CurrentIndex++
	renderGrammarQuestion()
}

func updateGrammarHUD() {
	doc := js.Global().Get("document")
	setElementText(doc, "grammarLiveScore", fmt.Sprintf("%d", gramState.Score))
	setElementText(doc, "grammarLiveStreak", fmt.Sprintf("%d 🔥", gramState.Streak))
	setElementText(doc, "grammarLiveMultiplier", fmt.Sprintf("%.1fx", gramState.Multiplier))

	acc := 0.0
	if gramState.AnsweredCount > 0 {
		acc = (float64(gramState.CorrectCount) / float64(gramState.AnsweredCount)) * 100.0
	}
	setElementText(doc, "grammarLiveAccuracy", fmt.Sprintf("%.0f%%", acc))
}

func finishGrammarRound() {
	gramState.Running = false
	doc := js.Global().Get("document")

	acc := 0.0
	if gramState.AnsweredCount > 0 {
		acc = (float64(gramState.CorrectCount) / float64(gramState.AnsweredCount)) * 100.0
	}

	setElementText(doc, "grammarPrompt", fmt.Sprintf("Round Completed! Final Score: %d Points • Accuracy: %.0f%%", gramState.Score, acc))

	optionsContainer := doc.Call("getElementById", "grammarOptionsContainer")
	if !optionsContainer.IsNull() {
		optionsContainer.Set("innerHTML", `<div style="padding: 24px; text-align: center; font-size: 1.05rem; color: #10b981; font-weight: 700;">Telemetry submitted to operator profile.</div>`)
	}

	btnNext := doc.Call("getElementById", "grammarNextBtn")
	if !btnNext.IsNull() {
		btnNext.Get("style").Set("display", "none")
	}

	SaveActivityResult(Result{
		ActivityType: "grammar",
		Score:        float64(gramState.Score),
		Accuracy:     acc,
		DurationMs:   int64(gramState.AnsweredCount * 6000),
		Timestamp:    time.Now().UnixMilli(),
	})
	triggerToast("Grammar Session Concluded", fmt.Sprintf("Final Score: %d Pts (Accuracy: %.0f%%)", gramState.Score, acc), "success")
}
