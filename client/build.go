package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type PageConfig struct {
	Title       string
	Description string
	Version     string
	Year        int
	Author      string
}

var defaultConfig = PageConfig{
	Title:       "SkillNODE",
	Description: "SkillNODE — Train hard. Look sharp. The premium skill platform for typing, coding, math, grammar, and multiplayer competition.",
	Version:     "1.0.1",
	Year:        time.Now().Year(),
	Author:      "Voxion Labs",
}

func main() {
	serveFlag := flag.Bool("serve", false, "Start a local development server")
	portFlag := flag.String("port", "5173", "Port for local dev server")
	wasmOnly := flag.Bool("wasm", false, "Only build the WebAssembly binary")
	tmplOnly := flag.Bool("tmpl", false, "Only build the HTML from Go templates")
	stylesOnly := flag.Bool("styles", false, "Only build CSS from Go templates")
	flag.Parse()

	if *serveFlag {
		buildAll()
		startServer(*portFlag)
		return
	}

	if *wasmOnly {
		buildWasm()
		ensureWasmExec()
		return
	}

	if *tmplOnly {
		buildTemplates()
		return
	}

	if *stylesOnly {
		buildStyles()
		return
	}

	buildAll()
}

func buildAll() {
	fmt.Println("==================================================")
	fmt.Println(" SkillNODE Go Build & Template Engine")
	fmt.Println("==================================================")

	buildTemplates()
	buildStyles()
	buildWasm()
	ensureWasmExec()
	syncDist()

	fmt.Println("==================================================")
	fmt.Println(" Build complete. Zero node_modules generated.")
	fmt.Println("==================================================")
}

func syncDist() {
	distDir := "dist"
	if err := os.MkdirAll(filepath.Join(distDir, "wasm"), 0755); err != nil {
		log.Printf("Warning creating dist directory: %v", err)
		return
	}

	copyFile("index.html", filepath.Join(distDir, "index.html"))
	copyFile("styles.css", filepath.Join(distDir, "styles.css"))
	copyFile(filepath.Join("wasm", "app.wasm"), filepath.Join(distDir, "wasm", "app.wasm"))
	copyFile(filepath.Join("wasm", "wasm_exec.js"), filepath.Join(distDir, "wasm", "wasm_exec.js"))
	if err := os.MkdirAll(filepath.Join(distDir, "static"), 0755); err == nil {
		copyFile(filepath.Join("static", "bridge.js"), filepath.Join(distDir, "static", "bridge.js"))
	}

	if _, err := os.Stat("logo.png"); err == nil {
		copyFile("logo.png", filepath.Join(distDir, "logo.png"))
	} else if _, err := os.Stat(filepath.Join("public", "logo.png")); err == nil {
		copyFile(filepath.Join("public", "logo.png"), filepath.Join(distDir, "logo.png"))
		copyFile(filepath.Join("public", "logo.png"), "logo.png")
	}

	fmt.Println("    Synced output to dist/ for Vercel deployment.")
}

func copyFile(src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		return
	}
	_ = os.WriteFile(dst, data, 0644)
}

func buildTemplates() {
	fmt.Println("--> Compiling Go templates (templates/*.tmpl -> index.html)...")

	// Collect HTML templates (exclude styles and wasm_exec tmpls)
	matches, err := filepath.Glob("templates/*.tmpl")
	if err != nil {
		log.Fatalf("Error scanning templates directory: %v", err)
	}

	var htmlTmpls []string
	for _, match := range matches {
		base := filepath.Base(match)
		if base == "styles.tmpl" || base == "wasm_exec.tmpl" {
			continue
		}
		htmlTmpls = append(htmlTmpls, match)
	}

	if len(htmlTmpls) == 0 {
		log.Fatalf("No HTML template files found in templates/")
	}

	t, err := template.ParseFiles(htmlTmpls...)
	if err != nil {
		log.Fatalf("Error parsing HTML templates: %v", err)
	}

	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf, "layout", defaultConfig)
	if err != nil {
		log.Fatalf("Error executing layout template: %v", err)
	}

	cleanHTML := compressHTML(buf.String())

	err = os.WriteFile("index.html", []byte(cleanHTML), 0644)
	if err != nil {
		log.Fatalf("Error writing index.html: %v", err)
	}

	fmt.Printf("    Generated index.html (%d bytes)\n", len(cleanHTML))
}

func buildStyles() {
	fmt.Println("--> Compiling CSS stylesheets from Go templates...")

	stylesContent, err := os.ReadFile("templates/styles.tmpl")
	if err != nil {
		log.Fatalf("Error reading templates/styles.tmpl: %v", err)
	}

	err = os.WriteFile("styles.css", stylesContent, 0644)
	if err != nil {
		log.Fatalf("Error writing styles.css: %v", err)
	}

	fmt.Printf("    Generated styles.css (%d bytes)\n", len(stylesContent))
}

func compressHTML(input string) string {
	reBlankLines := regexp.MustCompile(`\n\s*\n\s*\n+`)
	cleaned := reBlankLines.ReplaceAllString(input, "\n\n")
	return strings.TrimSpace(cleaned) + "\n"
}

func buildWasm() {
	fmt.Println("--> Compiling Go WebAssembly engine (src/ -> wasm/app.wasm)...")

	if err := os.MkdirAll("wasm", 0755); err != nil {
		log.Fatalf("Error creating wasm directory: %v", err)
	}

	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", "wasm/app.wasm", "./src")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error compiling WebAssembly: %v", err)
	}

	fi, err := os.Stat("wasm/app.wasm")
	if err != nil {
		log.Fatalf("Error reading wasm output file: %v", err)
	}

	duration := time.Since(start)
	fmt.Printf("    Compiled wasm/app.wasm: %.2f MB in %v\n", float64(fi.Size())/(1024*1024), duration.Round(time.Millisecond))
}

func ensureWasmExec() {
	if err := os.MkdirAll("wasm", 0755); err != nil {
		log.Fatalf("Error creating wasm directory: %v", err)
	}

	target := filepath.Join("wasm", "wasm_exec.js")

	// First try to copy from local template
	if content, err := os.ReadFile("templates/wasm_exec.tmpl"); err == nil {
		if err := os.WriteFile(target, content, 0644); err == nil {
			fmt.Println("    Provisioned wasm/wasm_exec.js from templates/wasm_exec.tmpl")
			return
		}
	}

	// Fallback to Go installation GOROOT
	cmd := exec.Command("go", "env", "GOROOT")
	out, err := cmd.Output()
	if err == nil {
		goroot := strings.TrimSpace(string(out))
		srcPath := filepath.Join(goroot, "lib", "wasm", "wasm_exec.js")
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			srcPath = filepath.Join(goroot, "misc", "wasm", "wasm_exec.js")
		}

		if content, err := os.ReadFile(srcPath); err == nil {
			if err := os.WriteFile(target, content, 0644); err == nil {
				fmt.Println("    Provisioned wasm/wasm_exec.js from GOROOT")
				return
			}
		}
	}

	fmt.Println("    Warning: Could not provision wasm_exec.js")
}

func startServer(port string) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Proper MIME type for wasm
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		} else if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		// Disable caching in dev mode
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fileServer.ServeHTTP(w, r)
	})

	addr := ":" + port
	fmt.Printf("\n🚀 Dev server running at http://localhost:%s/\n", port)
	fmt.Println("Press Ctrl+C to stop.")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
