//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("SkillNODE Go WebAssembly Engine Initializing...")

	// 1. Initialize local-first storage vault
	initStore()

	// 2. Initialize UI router, events, theme & profile
	initUI()

	// 3. Dismiss loading screen once UI is ready
	dismissLoadingScreen()

	// 4. Check Render backend API connectivity (non-blocking)
	CheckBackendHealth()

	fmt.Println("SkillNODE Go WebAssembly Engine Ready.")

	// Keep the Go runtime active to handle browser callbacks & events
	select {}
}

func dismissLoadingScreen() {
	fn := js.Global().Get("dismissSkillNodeLoader")
	if !fn.IsUndefined() && !fn.IsNull() {
		fn.Invoke()
		return
	}
	doc := js.Global().Get("document")
	loader := doc.Call("getElementById", "loadingScreen")
	if !loader.IsNull() {
		loader.Get("style").Set("opacity", "0")
		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			loader.Get("style").Set("display", "none")
			loader.Get("style").Set("visibility", "hidden")
			return nil
		}), 250)
	}
}
