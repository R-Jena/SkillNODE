//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func SyncResultToBackend(r Result) {
	apiURL := getAPIBaseURL() + "/api/results"
	p := GetSessionProfile()

	payload := map[string]interface{}{
		"activityType": r.ActivityType,
		"score":        r.Score,
		"accuracy":     r.Accuracy,
		"durationMs":   r.DurationMs,
		"userId":       p.ID,
		"username":     p.Username,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Non-blocking browser fetch
	go func() {
		opts := js.Global().Get("Object").New()
		opts.Set("method", "POST")

		headers := js.Global().Get("Object").New()
		headers.Set("Content-Type", "application/json")
		opts.Set("headers", headers)
		opts.Set("body", string(data))

		js.Global().Call("fetch", apiURL, opts)
	}()
}

func CheckBackendHealth() {
	apiURL := getAPIBaseURL() + "/healthz"

	go func() {
		promise := js.Global().Call("fetch", apiURL)
		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 && args[0].Get("ok").Bool() {
				fmt.Println("SkillNODE Render API Status: Online & Healthy")
			}
			return nil
		}))
	}()
}
