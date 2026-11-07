//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
)

func main() {
	fmt.Println("SkillNODE Go WebAssembly Engine Initializing...")

	// 1. Initialize local-first storage vault
	initStore()

	// 2. Initialize UI router, events, theme & profile
	initUI()

	fmt.Println("SkillNODE Go WebAssembly Engine Ready.")

	// Keep the Go runtime active to handle browser callbacks & events
	select {}
}
