//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"math"
	"strings"
	"syscall/js"
	"time"
)

type ComprehensionQuestion struct {
	Prompt      string
	Options     []string
	AnswerIndex int
	Rationale   string
}

type ReadingPassage struct {
	ID        string
	Title     string
	Category  string
	WordCount int
	Content   string
	Questions []ComprehensionQuestion
}

var readingBank = []ReadingPassage{
	{
		ID:        "raft-consensus",
		Title:     "Raft Consensus & State Machine Replication",
		Category:  "Distributed Systems",
		WordCount: 290,
		Content: `Distributed consensus algorithms ensure that a cluster of computing nodes agrees on a deterministic sequence of state machine transitions, even in the presence of network partitions or node failures. Raft was explicitly designed to address the conceptual opacity of Paxos by decomposing consensus into three independent sub-problems: leader election, log replication, and safety invariants.

In Raft, at any given moment, each node resides in one of three states: Leader, Follower, or Candidate. Under normal operations, a single elected leader coordinates all log append entries. Followers are passive; they issue no requests independently but respond to incoming remote procedure calls (RPCs) from leaders and candidates. 

A term acts as a logical clock, dividing execution into distinct chronological epochs identified by monotonically increasing integers. When a follower detects election timeout due to the absence of heartbeats, it transitions to the candidate state, increments the current term, votes for itself, and broadcasts RequestVote RPCs. If the candidate achieves a strict majority quorum, it becomes the new leader.

Safety is enforced through the Log Matching Property: if two logs contain an entry with the same index and term, they are identical up to that index. Furthermore, an elected leader is guaranteed to contain all committed entries from prior terms, preventing uncommitted candidate logs from overwriting historical truths.`,
		Questions: []ComprehensionQuestion{
			{
				Prompt:      "According to the passage, how does Raft prevent uncommitted candidate logs from overwriting historical cluster commits?",
				Options:     []string{
					"By requiring candidates to possess all committed entries from prior terms before being elected.",
					"By delegating log reconciliation to a secondary quorum outside the Raft cluster.",
					"By resetting all follower log indexes to zero whenever a leader timeout occurs.",
					"By utilizing physical wall-clock timestamps to order competing log entries.",
				},
				AnswerIndex: 0,
				Rationale:   "The safety invariant ensures that an elected leader is guaranteed to contain all committed entries from prior terms, preventing historical overwrite.",
			},
			{
				Prompt:      "What triggers a follower node to transition into the candidate state?",
				Options:     []string{
					"A direct instruction broadcast by the majority quorum.",
					"An election timeout caused by the absence of periodic leader heartbeats.",
					"The arrival of a client write request directly at the follower socket.",
					"A memory compaction warning issued by the underlying state machine.",
				},
				AnswerIndex: 1,
				Rationale:   "Followers transition to candidate state when an election timeout occurs due to missing leader heartbeats.",
			},
			{
				Prompt:      "What primary role do 'terms' serve within the Raft consensus model?",
				Options:     []string{
					"They designate physical disk write blocks on persistent storage drives.",
					"They serve as logical clocks dividing execution into monotonic chronological epochs.",
					"They define network subnet boundaries for encrypted peer-to-peer gossip.",
					"They count the total number of connected client HTTP sessions.",
				},
				AnswerIndex: 1,
				Rationale:   "Terms act as logical clocks dividing time into monotonic chronological epochs to detect obsolete leaders and stale candidates.",
			},
		},
	},
	{
		ID:        "wasm-stack-machine",
		Title:     "WebAssembly: Sandboxed Stack-Machine Execution",
		Category:  "Compiler & Runtime Design",
		WordCount: 310,
		Content: `WebAssembly (WASM) represents a binary instruction format designed as a portable, size- and time-efficient compilation target for high-performance applications executing in sandboxed environments. Unlike register-based virtual machines or JavaScript's dynamic AST interpretation, WASM operates upon a structured stack-machine abstraction.

In this architecture, instructions push values onto an evaluation stack, perform arithmetic or logical operations upon those operands, and pop intermediate results back into memory or destination stack frames. This stack-based representation substantially compacts bytecode encoding and enables rapid linear-time single-pass validation and compilation into native machine instructions via just-in-time (JIT) or ahead-of-time (AOT) compilers.

A foundational security guarantee of WebAssembly is memory isolation. A WASM module cannot access arbitrary host process memory. Instead, all load and store instructions operate strictly within an untyped, linearly addressed array of contiguous bytes known as Linear Memory. 

This memory buffer is instantiated by the host environment and can be dynamically grown in 64-kilobyte page increments. Every memory access is subjected to hardware- or software-enforced boundary checks. Any attempt to dereference addresses exceeding the allocated linear bounds triggers an immediate uncatchable runtime trap, isolating untrusted guest code from host memory.`,
		Questions: []ComprehensionQuestion{
			{
				Prompt:      "How does WebAssembly enforce strict memory isolation against host processes?",
				Options:     []string{
					"By restricting all memory operations to a sandboxed Linear Memory buffer with boundary checks.",
					"By running all guest code in separate operating system hypervisor virtual machines.",
					"By forbidding pointers and references entirely from the bytecode specification.",
					"By encrypting all stack frames with asymmetric public-key cryptography.",
				},
				AnswerIndex: 0,
				Rationale:   "WASM operates exclusively within an untyped, linearly addressed memory buffer where every access undergoes strict boundary checks.",
			},
			{
				Prompt:      "What is an architectural benefit of WebAssembly's stack-machine representation?",
				Options:     []string{
					"It eliminates the need for any compiling or native code generation.",
					"It allows compact bytecode encoding and linear-time single-pass validation.",
					"It automatically optimizes SQL queries before sending them to the host.",
					"It replaces garbage collection in high-level managed languages completely.",
				},
				AnswerIndex: 1,
				Rationale:   "Stack-based bytecode enables compact encoding and fast linear-time single-pass compilation into native machine code.",
			},
			{
				Prompt:      "In what unit increments can WebAssembly Linear Memory be dynamically expanded?",
				Options:     []string{"64-kilobyte pages", "4-kilobyte blocks", "1-megabyte chunks", "512-byte sectors"},
				AnswerIndex: 0,
				Rationale:   "The WASM specification defines memory growth in standard 64 KB page increments.",
			},
		},
	},
	{
		ID:        "zero-knowledge-snarks",
		Title:     "Zero-Knowledge Cryptography & zk-SNARKs",
		Category:  "Applied Cryptography",
		WordCount: 305,
		Content: `Zero-knowledge cryptography enables a prover to mathematically convince a verifier that a computational statement is true without disclosing any secret witness data underpinning that assertion. Zero-Knowledge Succinct Non-Interactive Arguments of Knowledge (zk-SNARKs) represent the pinnacle of this paradigm, offering proof sizes of several hundred bytes and verification times in single-digit milliseconds.

To construct a zk-SNARK proof, a high-level computational program is first compiled into an algebraic intermediate representation—typically Rank-1 Constraint Systems (R1CS) or Arithmetic Circuits over finite fields. In an arithmetic circuit, addition and multiplication gates represent individual computation steps. 

These constraints are subsequently transformed via polynomial interpolation into a Quadratic Arithmetic Program (QAP). The prover demonstrates knowledge of valid polynomials that satisfy all circuit constraints by evaluating them at a secret evaluation point concealed through elliptic curve pairings or homomorphic polynomial commitments.

The property of 'succinctness' requires that the proof size remains constant or logarithmic relative to computation size, and verification time must be negligible compared to re-executing the raw computation. Consequently, zk-SNARKs serve as the cornerstone of blockchain layer-2 rollups and verifiable cloud computing offloading.`,
		Questions: []ComprehensionQuestion{
			{
				Prompt:      "What does 'succinctness' signify in the context of zk-SNARK proofs?",
				Options:     []string{
					"The proof size is small (constant/logarithmic) and verification time is negligible.",
					"The cryptographic key requires zero initial generation setup.",
					"The underlying algorithm only operates on floating-point arithmetic numbers.",
					"The prover and verifier must engage in continuous multi-turn round-trip communication.",
				},
				AnswerIndex: 0,
				Rationale:   "Succinctness dictates that proof size is minimal (hundreds of bytes) and verification is orders of magnitude faster than re-running the computation.",
			},
			{
				Prompt:      "Which mathematical structure represents computational steps before transformation into Quadratic Arithmetic Programs?",
				Options:     []string{
					"Arithmetic circuits and Rank-1 Constraint Systems (R1CS)",
					"Directed Acyclic Graph traversal trees in binary format",
					"Merkle-Patricia radix tries with keccak256 hashes",
					"Relational B-Tree database indexes with foreign keys",
				},
				AnswerIndex: 0,
				Rationale:   "Programs are compiled into Arithmetic Circuits or Rank-1 Constraint Systems (R1CS) over finite fields before QAP transformation.",
			},
			{
				Prompt:      "What is the primary objective of a zero-knowledge proof?",
				Options:     []string{
					"To compress large video files without loss of visual quality.",
					"To prove a statement's truth without revealing the underlying private witness data.",
					"To execute database transactions across multi-region data centers instantly.",
					"To generate decentralized public keys from physical biometric scans.",
				},
				AnswerIndex: 1,
				Rationale:   "Zero-knowledge proofs allow a prover to establish the veracity of a statement without exposing secret witness data.",
			},
		},
	},
	{
		ID:        "hardware-cache-coherence",
		Title:     "Hardware Cache Hierarchies & False Sharing",
		Category:  "Computer Architecture",
		WordCount: 320,
		Content: `Modern multi-core processors mitigate the widening latency disparity between CPU execution units and dynamic RAM through a hierarchical memory subsystem consisting of L1, L2, and shared L3 caches. L1 caches offer access latencies of approximately 1 nanosecond, while main memory fetches frequently incur latencies exceeding 60 to 100 nanoseconds.

To maximize throughput, processors do not transfer individual bytes; memory transactions occur in discrete units termed Cache Lines—typically 64 bytes in modern x86-64 and ARM architectures. Cache coherence protocols (such as MESI or MOESI) maintain memory consistency across competing processor cores by coordinating the state of each line (Modified, Exclusive, Shared, Invalid).

A prevalent performance hazard in concurrent programming is 'False Sharing'. False sharing arises when two independent threads executing on separate CPU cores read and modify completely unrelated variables that happen to reside within the identical 64-byte cache line. 

When Core A writes to variable X, the coherence controller invalidates the entire cache line across all other cores. When Core B attempts to read or write neighboring variable Y on that same line, it experiences an unexpected cache miss, triggering an expensive bus synchronization round-trip. Mitigating false sharing necessitates explicit hardware padding or alignment attributes to ensure concurrent variables inhabit separate cache boundaries.`,
		Questions: []ComprehensionQuestion{
			{
				Prompt:      "What underlying condition produces the performance anomaly known as 'False Sharing'?",
				Options:     []string{
					"Independent threads modify unrelated variables that reside within the same 64-byte cache line.",
					"Two cores attempt to read from write-only flash memory sectors simultaneously.",
					"An operating system scheduler assigns two processes identical virtual memory addresses.",
					"A single-threaded program runs out of stack space and overflows into the heap.",
				},
				AnswerIndex: 0,
				Rationale:   "False sharing occurs when unrelated variables accessed by distinct cores share the same cache line, triggering unnecessary coherence invalidations.",
			},
			{
				Prompt:      "What is the standard unit size of memory transferred between CPU caches and main memory in modern architectures?",
				Options:     []string{"64 bytes", "1 byte", "4 kilobytes", "512 bytes"},
				AnswerIndex: 0,
				Rationale:   "Modern architectures transfer memory in standard 64-byte cache lines.",
			},
			{
				Prompt:      "How can software engineers proactively prevent false sharing between concurrent variables?",
				Options:     []string{
					"By applying memory padding or cache-line alignment to separate variables onto distinct lines.",
					"By removing all L1 and L2 caches and querying dynamic RAM directly.",
					"By converting all multi-threaded operations into single-threaded asynchronous promises.",
					"By increasing CPU clock speeds beyond thermal throttling limits.",
				},
				AnswerIndex: 0,
				Rationale:   "Explicit memory padding or cache-line alignment ensures concurrent variables occupy distinct cache lines.",
			},
		},
	},
	{
		ID:        "async-event-loops",
		Title:     "Asynchronous Event Loops vs Kernel Threading",
		Category:  "Systems Programming",
		WordCount: 315,
		Content: `High-concurrency network servers historically relied upon multithreaded architectures, assigning an operating system thread to each incoming TCP connection. While conceptually straightforward, thread-per-connection architectures scale poorly under tens of thousands of concurrent sockets due to substantial kernel memory footprints (often 1 to 8 megabytes of stack per thread) and pervasive context-switching overhead.

In contrast, non-blocking asynchronous event loops leverage OS-level multiplexing primitives—such as epoll on Linux, kqueue on BSD/macOS, and IOCP on Windows. Rather than blocking an OS thread during network I/O, the event loop registers file descriptors with the kernel notification subsystem and continues processing ready tasks.

When an I/O event completes, the kernel alerts the event loop, which dispatches the associated callback or resumes the suspended coroutine on a lightweight single thread or small worker pool. This paradigm eliminates preemptive context switches and eliminates thread race conditions within the main execution thread.

However, event loops introduce cooperative scheduling vulnerabilities: if an asynchronous callback performs CPU-intensive computation or enters an unyielding loop, it blocks the entire event loop, starving all concurrent I/O handlers. Hybrid systems frequently combine event loops with thread worker pools to offload intensive tasks.`,
		Questions: []ComprehensionQuestion{
			{
				Prompt:      "What is a primary limitation of the thread-per-connection architecture when scaling to tens of thousands of concurrent clients?",
				Options:     []string{
					"High memory consumption per thread stack and kernel context-switching overhead.",
					"The inability of operating systems to allocate more than 256 TCP port numbers.",
					"Incompatibility with modern solid-state NVMe storage controllers.",
					"Mandatory hardware cryptographic signing requirements on each socket.",
				},
				AnswerIndex: 0,
				Rationale:   "Thread-per-connection architectures suffer from memory stack bloat and high kernel context-switching costs at scale.",
			},
			{
				Prompt:      "What critical vulnerability must developers avoid when writing code inside an asynchronous event loop?",
				Options:     []string{
					"Executing long CPU-intensive computations that monopolize the loop thread.",
					"Using non-blocking network socket reads instead of blocking system calls.",
					"Registering multiple file descriptors with the kernel epoll multiplexer.",
					"Employing small worker thread pools for background disk access.",
				},
				AnswerIndex: 0,
				Rationale:   "Because cooperative event loops run on shared threads, CPU-intensive unyielding work halts the entire loop and starves other handlers.",
			},
			{
				Prompt:      "Which OS-level multiplexing primitive powers high-performance asynchronous event loops on Linux?",
				Options:     []string{"epoll", "kqueue", "IOCP", "select16"},
				AnswerIndex: 0,
				Rationale:   "Linux utilizes epoll as its scalable I/O event notification multiplexing facility.",
			},
		},
	},
	{
		ID:        "transformers-attention",
		Title:     "The Architecture of Scalable Self-Attention",
		Category:  "Artificial Intelligence",
		WordCount: 325,
		Content: `The Transformer architecture supplanted Recurrent Neural Networks (RNNs) as the dominant foundation for natural language processing and multimodal foundation models by discarding sequential recurrence entirely in favor of Scaled Dot-Product Self-Attention. 

In sequential models, processing token at step t inherently depended on the hidden state of step t-1, creating an algorithmic bottleneck that precluded parallel computation across GPU streaming multiprocessors. 

In self-attention, input token embeddings are linearly projected into three distinct vector representations: Queries (Q), Keys (K), and Values (V). The attention weights are computed by calculating the dot product of all queries with all keys, scaling the resulting matrix by the inverse square root of the key dimension (√d_k) to prevent vanishing gradients in softmax regions, and applying the softmax function to derive a probability distribution over all sequence positions.

This formulation allows every token in an input sequence to attend directly to every other token in O(1) sequential operations, enabling massive parallelization across tens of thousands of accelerator chips during training. 

However, standard self-attention exhibits quadratic O(N²) time and memory complexity with respect to sequence length N, prompting modern innovations such as FlashAttention—which leverages SRAM tiling to minimize high-bandwidth memory (HBM) read-write rounds—and linear attention approximations.`,
		Questions: []ComprehensionQuestion{
			{
				Prompt:      "Why did the self-attention mechanism enable far greater training parallelization than Recurrent Neural Networks?",
				Options:     []string{
					"It eliminated sequential step-by-step state dependencies, allowing parallel matrix operations across tokens.",
					"It reduced all neural network weights to 1-bit integers.",
					"It bypassed the need for gradient descent backpropagation entirely.",
					"It replaced GPU tensor cores with standard central processing units.",
				},
				AnswerIndex: 0,
				Rationale:   "Unlike sequential RNNs that require token t-1 before computing t, self-attention processes all tokens concurrently via matrix multiplications.",
			},
			{
				Prompt:      "What is the theoretical algorithmic time and memory complexity of standard full self-attention with sequence length N?",
				Options:     []string{"Quadratic O(N²)", "Linear O(N)", "Logarithmic O(log N)", "Constant O(1)"},
				AnswerIndex: 0,
				Rationale:   "Comparing every token's Query with all Keys creates an N×N attention matrix, resulting in quadratic O(N²) complexity.",
			},
			{
				Prompt:      "Why is the dot product of Queries and Keys scaled by 1/√d_k in the attention equation?",
				Options:     []string{
					"To prevent large magnitude dot products from pushing softmax into regions with vanishing gradients.",
					"To convert negative vectors into positive real values before matrix inversion.",
					"To encrypt token embeddings before transferring them over PCI Express channels.",
					"To force all attention matrix diagonal values to zero.",
				},
				AnswerIndex: 0,
				Rationale:   "Scaling by 1/√d_k prevents large values from driving the softmax function into extremely small gradient regions.",
			},
		},
	},
}

type ReadingState struct {
	ActivePassageIndex int
	ReadingStartTime   int64
	ReadingFinished    bool
	ReadWPM            int
	Answers            map[int]int
	Submitted          bool
	Score              int
}

var rState = &ReadingState{
	ActivePassageIndex: 0,
	Answers:            make(map[int]int),
}
var readingInitialized = false

func initReadingView() {
	if readingInitialized {
		renderActivePassage()
		return
	}
	readingInitialized = true

	doc := js.Global().Get("document")

	// Done Reading Button
	doneBtn := doc.Call("getElementById", "readingDoneBtn")
	if !doneBtn.IsNull() {
		doneBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			finishReadingPhase()
			return nil
		}))
	}

	// Submit Answers Button
	submitBtn := doc.Call("getElementById", "readingSubmitBtn")
	if !submitBtn.IsNull() {
		submitBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			submitReadingAnswers()
			return nil
		}))
	}

	// Reset / Next Passage Button
	nextPassageBtn := doc.Call("getElementById", "readingNextPassageBtn")
	if !nextPassageBtn.IsNull() {
		nextPassageBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			rState.ActivePassageIndex = (rState.ActivePassageIndex + 1) % len(readingBank)
			renderActivePassage()
			return nil
		}))
	}

	renderPassageTabs()
	renderActivePassage()
}

func renderPassageTabs() {
	doc := js.Global().Get("document")
	tabsContainer := doc.Call("getElementById", "readingPassageTabs")
	if tabsContainer.IsNull() {
		return
	}

	var sb strings.Builder
	for i, p := range readingBank {
		activeStyle := "background: var(--surface-2); color: var(--text-muted); border: 1px solid var(--border-subtle);"
		if i == rState.ActivePassageIndex {
			activeStyle = "background: var(--accent-dim); color: var(--text); border: 1px solid var(--accent); font-weight: 700;"
		}
		sb.WriteString(fmt.Sprintf(`
			<button class="btn btn-sm reading-tab-btn" data-passage-idx="%d" style="border-radius: 999px; height: 32px; font-size: 0.78rem; %s gap: 6px; flex-shrink: 0;">
				<span>%s</span>
				<span style="font-size: 0.68rem; color: #a3e635;">• %d words</span>
			</button>
		`, i, activeStyle, p.Title, p.WordCount))
	}
	tabsContainer.Set("innerHTML", sb.String())

	// Attach click listeners
	buttons := tabsContainer.Call("querySelectorAll", "[data-passage-idx]")
	bLen := buttons.Get("length").Int()
	for i := 0; i < bLen; i++ {
		btn := buttons.Index(i)
		btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			idx := this.Call("getAttribute", "data-passage-idx").Int()
			rState.ActivePassageIndex = idx
			renderActivePassage()
			renderPassageTabs()
			return nil
		}))
	}
}

func renderActivePassage() {
	p := readingBank[rState.ActivePassageIndex]
	doc := js.Global().Get("document")

	rState.ReadingStartTime = time.Now().UnixMilli()
	rState.ReadingFinished = false
	rState.ReadWPM = 0
	rState.Answers = make(map[int]int)
	rState.Submitted = false
	rState.Score = 0

	setElementText(doc, "readingCategoryBadge", p.Category)
	setElementText(doc, "readingTitleDisplay", p.Title)
	setElementText(doc, "readingWordCountBadge", fmt.Sprintf("%d words", p.WordCount))
	setElementText(doc, "readingPassageBody", p.Content)

	// Reset HUD
	setElementText(doc, "readingWpmDisplay", "Reading...")
	setElementText(doc, "readingScoreDisplay", "0 pts")
	setElementText(doc, "readingAccuracyDisplay", "-")

	// Hide Questions container until user finishes reading
	qContainer := doc.Call("getElementById", "readingQuestionsContainer")
	if !qContainer.IsNull() {
		qContainer.Get("style").Set("display", "none")
	}

	doneBtn := doc.Call("getElementById", "readingDoneBtn")
	if !doneBtn.IsNull() {
		doneBtn.Get("style").Set("display", "inline-flex")
		doneBtn.Set("disabled", false)
	}

	submitBtn := doc.Call("getElementById", "readingSubmitBtn")
	if !submitBtn.IsNull() {
		submitBtn.Get("style").Set("display", "none")
	}

	feedbackEl := doc.Call("getElementById", "readingFeedbackFeed")
	if !feedbackEl.IsNull() {
		feedbackEl.Set("innerHTML", "")
	}
}

func finishReadingPhase() {
	if rState.ReadingFinished {
		return
	}
	rState.ReadingFinished = true

	now := time.Now().UnixMilli()
	elapsedSeconds := float64(now-rState.ReadingStartTime) / 1000.0
	if elapsedSeconds < 1.0 {
		elapsedSeconds = 1.0
	}

	p := readingBank[rState.ActivePassageIndex]
	wpm := int((float64(p.WordCount) / elapsedSeconds) * 60.0)
	if wpm > 1200 {
		wpm = 1200
	}
	rState.ReadWPM = wpm

	doc := js.Global().Get("document")
	setElementText(doc, "readingWpmDisplay", fmt.Sprintf("%d WPM", wpm))

	doneBtn := doc.Call("getElementById", "readingDoneBtn")
	if !doneBtn.IsNull() {
		doneBtn.Get("style").Set("display", "none")
	}

	// Show Questions
	qContainer := doc.Call("getElementById", "readingQuestionsContainer")
	if !qContainer.IsNull() {
		qContainer.Get("style").Set("display", "flex")
		renderComprehensionQuestions()
	}

	submitBtn := doc.Call("getElementById", "readingSubmitBtn")
	if !submitBtn.IsNull() {
		submitBtn.Get("style").Set("display", "inline-flex")
	}

	triggerToast("Reading Pace Logged", fmt.Sprintf("Calculated %d WPM. Answer the comprehension questions.", wpm), "info")
}

func renderComprehensionQuestions() {
	p := readingBank[rState.ActivePassageIndex]
	doc := js.Global().Get("document")
	feed := doc.Call("getElementById", "readingQuestionsFeed")
	if feed.IsNull() {
		return
	}

	var sb strings.Builder
	for qIdx, q := range p.Questions {
		sb.WriteString(fmt.Sprintf(`
			<div class="card-inner" style="padding: 20px; display: flex; flex-direction: column; gap: 14px;">
				<div style="font-weight: 700; font-size: 0.95rem; color: var(--text);">
					<span style="color: var(--accent); margin-right: 6px;">Q%d.</span> %s
				</div>
				<div style="display: flex; flex-direction: column; gap: 8px;">
		`, qIdx+1, q.Prompt))

		for optIdx, opt := range q.Options {
			sb.WriteString(fmt.Sprintf(`
				<label style="display: flex; align-items: flex-start; gap: 10px; cursor: pointer; padding: 10px 14px; border-radius: 10px; background: var(--surface); border: 1px solid var(--border-subtle); transition: border-color 0.15s;">
					<input type="radio" name="comprehension_q_%d" value="%d" style="margin-top: 3px;" />
					<span style="font-size: 0.88rem; color: var(--text); line-height: 1.5;">%s</span>
				</label>
			`, qIdx, optIdx, opt))
		}

		sb.WriteString(`
				</div>
			</div>
		`)
	}

	feed.Set("innerHTML", sb.String())
}

func submitReadingAnswers() {
	if rState.Submitted {
		return
	}

	p := readingBank[rState.ActivePassageIndex]
	doc := js.Global().Get("document")

	correctCount := 0
	totalQuestions := len(p.Questions)

	for qIdx, q := range p.Questions {
		radios := doc.Call("querySelectorAll", fmt.Sprintf("input[name='comprehension_q_%d']", qIdx))
		rLen := radios.Get("length").Int()
		chosen := -1
		for i := 0; i < rLen; i++ {
			r := radios.Index(i)
			if r.Get("checked").Bool() {
				chosen = i
				break
			}
		}

		if chosen == -1 {
			triggerToast("Incomplete Submission", fmt.Sprintf("Please answer Question %d before submitting.", qIdx+1), "warning")
			return
		}

		rState.Answers[qIdx] = chosen
		if chosen == q.AnswerIndex {
			correctCount++
		}
	}

	rState.Submitted = true
	accuracy := (float64(correctCount) / float64(totalQuestions)) * 100.0
	scorePts := int(accuracy * 10.0)
	rState.Score = scorePts

	setElementText(doc, "readingScoreDisplay", fmt.Sprintf("%d pts", scorePts))
	setElementText(doc, "readingAccuracyDisplay", fmt.Sprintf("%.0f%%", accuracy))

	// Show Answer Explanations Feed
	feedbackEl := doc.Call("getElementById", "readingFeedbackFeed")
	if !feedbackEl.IsNull() {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf(`
			<div style="padding: 16px; border-radius: 14px; background: rgba(16, 185, 129, 0.12); border: 1px solid rgba(16, 185, 129, 0.25); margin-top: 14px;">
				<div style="font-size: 1rem; font-weight: 800; color: #10b981;">
					Comprehension Evaluation: %d/%d Correct (%.0f%%) • Velocity: %d WPM
				</div>
				<div style="font-size: 0.82rem; color: var(--text-muted); margin-top: 4px;">
					Performance telemetry recorded and stored in local HMAC cryptographic envelope.
				</div>
			</div>
		`, correctCount, totalQuestions, accuracy, rState.ReadWPM))

		for qIdx, q := range p.Questions {
			chosen := rState.Answers[qIdx]
			pass := (chosen == q.AnswerIndex)
			color := "#10b981"
			badge := "PASSED"
			if !pass {
				color = "#ef4444"
				badge = "MISSED"
			}
			sb.WriteString(fmt.Sprintf(`
				<div style="padding: 12px 14px; border-radius: 10px; background: var(--surface-2); border: 1px solid var(--border-subtle); margin-top: 10px;">
					<div style="display: flex; align-items: center; justify-content: space-between;">
						<strong style="color: var(--text); font-size: 0.88rem;">Question %d</strong>
						<span style="color: %s; font-weight: 700; font-size: 0.78rem;">%s</span>
					</div>
					<div style="font-size: 0.82rem; color: var(--text-muted); margin-top: 4px; line-height: 1.5;">
						<strong>Rationale:</strong> %s
					</div>
				</div>
			`, qIdx+1, color, badge, q.Rationale))
		}
		feedbackEl.Set("innerHTML", sb.String())
	}

	submitBtn := doc.Call("getElementById", "readingSubmitBtn")
	if !submitBtn.IsNull() {
		submitBtn.Get("style").Set("display", "none")
	}

	// Save result to store
	SaveActivityResult(Result{
		ActivityType: "comprehension",
		Score:        float64(scorePts),
		Accuracy:     accuracy,
		DurationMs:   int64(math.Max(float64(p.WordCount*60000/int(math.Max(float64(rState.ReadWPM), 100))), 15000)),
		Timestamp:    time.Now().UnixMilli(),
	})

	triggerToast("Assessment Submitted 🎉", fmt.Sprintf("Score: %d Pts • Accuracy: %.0f%% • Rate: %d WPM", scorePts, accuracy, rState.ReadWPM), "success")
}
