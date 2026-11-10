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
		defer func() {
			_ = recover()
		}()

		opts := js.Global().Get("Object").New()
		opts.Set("method", "POST")

		headers := js.Global().Get("Object").New()
		headers.Set("Content-Type", "application/json")
		opts.Set("headers", headers)
		opts.Set("body", string(data))

		promise := js.Global().Call("fetch", apiURL, opts)
		promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}))
	}()
}

func CheckBackendHealth() {
	apiURL := getAPIBaseURL() + "/healthz"

	go func() {
		defer func() {
			_ = recover()
		}()

		promise := js.Global().Call("fetch", apiURL)
		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 && args[0].Get("ok").Bool() {
				fmt.Println("SkillNODE Render API Status: Online & Healthy")
			}
			return nil
		}))
		promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// Backend is asleep or standby - handled silently without throwing
			return nil
		}))
	}()
}
