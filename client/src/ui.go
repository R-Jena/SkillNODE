//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

var currentActiveRoute = "dashboard"

func initUI() {
	doc := js.Global().Get("document")

	// 1. Sync theme on startup
	applyStoredTheme()

	// 2. Sync profile in sidebar and profile view
	renderProfileInfo()

	// 3. Register route handler on window for the JS bridge
	js.Global().Set("onRouteChanged", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			SwitchView(args[0].String())
		}
		return nil
	}))

	// 4. Initial route from location.hash
	initialHash := js.Global().Get("location").Get("hash").String()
	cleanRoute := strings.TrimPrefix(initialHash, "#/")
	cleanRoute = strings.TrimPrefix(cleanRoute, "#")
	if cleanRoute == "" {
		cleanRoute = "dashboard"
	}
	SwitchView(cleanRoute)

	// 5. Theme toggle button listener
	themeBtn := doc.Call("getElementById", "toggleThemeBtn")
	if !themeBtn.IsNull() {
		themeBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			toggleTheme()
			return nil
		}))
	}

	// 6. Auth modal triggers
	authBtn := doc.Call("getElementById", "authTriggerBtn")
	if !authBtn.IsNull() {
		authBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			openModal("authModalBackdrop")
			return nil
		}))
	}

	authCloseBtn := doc.Call("getElementById", "authModalCloseBtn")
	if !authCloseBtn.IsNull() {
		authCloseBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			closeModal("authModalBackdrop")
			return nil
		}))
	}

	guestBtn := doc.Call("getElementById", "authGuestBtn")
	if !guestBtn.IsNull() {
		guestBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			p := createDefaultGuestProfile()
			saveProfileToStorage(p)
			renderProfileInfo()
			closeModal("authModalBackdrop")
			triggerToast("Session Initialized", "Guest operator profile active.", "success")
			return nil
		}))
	}

	// 7. Profile view edit modal trigger
	editProfileBtn := doc.Call("getElementById", "profileEditModalBtn")
	if !editProfileBtn.IsNull() {
		editProfileBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			openModal("authModalBackdrop")
			return nil
		}))
	}

	shareBtn := doc.Call("getElementById", "profileShareBtn")
	if !shareBtn.IsNull() {
		shareBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			p := GetSessionProfile()
			url := js.Global().Get("location").Get("origin").String() + "#/u/" + p.ID
			navigator := js.Global().Get("navigator")
			if !navigator.Get("clipboard").IsUndefined() {
				navigator.Get("clipboard").Call("writeText", url)
				triggerToast("Card Link Copied", "Public share URL copied to clipboard.", "success")
			} else {
				triggerToast("Operator Profile", url, "info")
			}
			return nil
		}))
	}

	// 8. Bind room presets in multiplayer view
	bindRoomPresetButtons()
}

func SwitchView(route string) {
	route = strings.TrimPrefix(route, "#/")
	route = strings.TrimPrefix(route, "#")
	route = strings.TrimPrefix(route, "/")
	if route == "" {
		route = "dashboard"
	}
	currentActiveRoute = route

	doc := js.Global().Get("document")

	// Hide all views
	views := doc.Call("querySelectorAll", ".page-view")
	vCount := views.Get("length").Int()
	for i := 0; i < vCount; i++ {
		views.Index(i).Get("style").Set("display", "none")
	}

	// Resolve target view element ID
	targetViewID := "view-dashboard"
	viewTitle := "Dashboard"
	navKey := "dashboard"

	switch {
	case route == "dashboard":
		targetViewID = "view-dashboard"
		viewTitle = "Dashboard & Telemetry"
		navKey = "dashboard"
		renderDashboardSummary()
	case route == "activities":
		targetViewID = "view-activities"
		viewTitle = "Skill Training Labs"
		navKey = "activities"
	case strings.HasPrefix(route, "activities/typing"):
		targetViewID = "view-activity-typing"
		viewTitle = "Speed Lab: Typing Velocity"
		navKey = "activities"
		initTypingView()
	case strings.HasPrefix(route, "activities/code"):
		targetViewID = "view-activity-code"
		viewTitle = "Developer Track: Code Arena"
		navKey = "activities"
	case strings.HasPrefix(route, "activities/math"):
		targetViewID = "view-activity-math"
		viewTitle = "Quant Lab: Math Clash"
		navKey = "activities"
	case strings.HasPrefix(route, "multiplayer"):
		targetViewID = "view-multiplayer"
		viewTitle = "Battle Net: Real-Time Multiplayer"
		navKey = "multiplayer"
	case strings.HasPrefix(route, "profile"):
		targetViewID = "view-profile"
		viewTitle = "Operator Profile Studio"
		navKey = "profile"
		renderProfileInfo()
	default:
		targetViewID = "view-dashboard"
		viewTitle = "Dashboard"
		navKey = "dashboard"
	}

	// Show target view
	target := doc.Call("getElementById", targetViewID)
	if !target.IsNull() {
		target.Get("style").Set("display", "flex")
	}

	// Update Title in top bar
	titleEl := doc.Call("getElementById", "currentViewTitle")
	if !titleEl.IsNull() {
		titleEl.Set("innerText", viewTitle)
	}

	// Update active class in sidebar nav
	navItems := doc.Call("querySelectorAll", ".shell-nav .nav-item")
	nCount := navItems.Get("length").Int()
	for i := 0; i < nCount; i++ {
		item := navItems.Index(i)
		itemKey := item.Call("getAttribute", "data-nav").String()
		if itemKey == navKey {
			item.Get("classList").Call("add", "active")
		} else {
			item.Get("classList").Call("remove", "active")
		}
	}

	// Scroll frame to top
	frame := doc.Call("getElementById", "pageFrame")
	if !frame.IsNull() {
		frame.Set("scrollTop", 0)
	}
}

func applyStoredTheme() {
	doc := js.Global().Get("document")
	theme := GetStoredTheme()
	doc.Get("documentElement").Call("setAttribute", "data-theme", theme)

	icon := doc.Call("getElementById", "themeIcon")
	if !icon.IsNull() {
		if strings.Contains(theme, "light") {
			icon.Set("className", "fas fa-sun")
		} else {
			icon.Set("className", "fas fa-moon")
		}
	}
}

func toggleTheme() {
	doc := js.Global().Get("document")
	current := doc.Get("documentElement").Call("getAttribute", "data-theme").String()

	newTheme := "dark-purple"
	if current == "dark" || current == "dark-purple" {
		newTheme = "light-white"
	}

	doc.Get("documentElement").Call("setAttribute", "data-theme", newTheme)
	SetStoredTheme(newTheme)

	icon := doc.Call("getElementById", "themeIcon")
	if !icon.IsNull() {
		if strings.Contains(newTheme, "light") {
			icon.Set("className", "fas fa-sun")
		} else {
			icon.Set("className", "fas fa-moon")
		}
	}

	triggerToast("Theme Updated", fmt.Sprintf("Active theme: %s", newTheme), "info")
}

func renderProfileInfo() {
	p := GetSessionProfile()
	if p == nil {
		return
	}

	doc := js.Global().Get("document")

	// Sidebar
	setElementText(doc, "sidebarPlayerName", p.Name)
	setElementText(doc, "sidebarPlayerHandle", "@"+p.Username)
	avatarEl := doc.Call("getElementById", "sidebarAvatar")
	if !avatarEl.IsNull() && len(p.Name) > 0 {
		avatarEl.Set("innerText", strings.ToUpper(string(p.Name[0])))
	}

	// Profile View
	setElementText(doc, "profileNameDisplay", p.Name)
	setElementText(doc, "profileHandleDisplay", "@"+p.Username)
	setElementText(doc, "profileBioDisplay", p.Bio)
	setElementText(doc, "profileSessionId", p.ID)
	setElementText(doc, "profileRegion", p.Region)

	profileRing := doc.Call("getElementById", "profileAvatarRing")
	if !profileRing.IsNull() && len(p.Name) > 0 {
		profileRing.Set("innerText", strings.ToUpper(string(p.Name[0])))
	}
}

func renderDashboardSummary() {
	results := LoadResults()
	doc := js.Global().Get("document")

	// Total attempts
	setElementText(doc, "metricTotalAttempts", fmt.Sprintf("%d", len(results)))

	// Calculate best WPM and accuracy
	bestWpm := 0.0
	for _, r := range results {
		if r.ActivityType == "typing" && r.Score > bestWpm {
			bestWpm = r.Score
		}
	}

	if bestWpm > 0 {
		setElementText(doc, "profileBestWpm", fmt.Sprintf("%.0f WPM", bestWpm))
		setElementText(doc, "metricPersonalBests", fmt.Sprintf("%.0f WPM", bestWpm))
	}

	// Populate Recent List
	recentContainer := doc.Call("getElementById", "dashboardRecentFeed")
	if !recentContainer.IsNull() && len(results) > 0 {
		var html strings.Builder
		limit := len(results)
		if limit > 6 {
			limit = 6
		}
		for i := 0; i < limit; i++ {
			r := results[i]
			html.WriteString(fmt.Sprintf(`
				<div style="display: flex; align-items: center; justify-content: space-between; padding: 10px 14px; border-radius: 12px; background: var(--surface-2); border: 1px solid var(--border-subtle);">
					<div>
						<div style="font-size: 0.88rem; font-weight: 700; color: var(--text); text-transform: capitalize;">%s</div>
						<div style="font-size: 0.72rem; color: var(--text-faint);">Score: %.0f • Acc: %.0f%%</div>
					</div>
					<span class="badge" style="background: rgba(139, 92, 246, 0.15); color: var(--accent-bright);">%.1fs</span>
				</div>
			`, r.ActivityType, r.Score, r.Accuracy, float64(r.DurationMs)/1000.0))
		}
		recentContainer.Set("innerHTML", html.String())
	}
}

func bindRoomPresetButtons() {
	doc := js.Global().Get("document")
	presets := doc.Call("querySelectorAll", "[data-room-preset]")
	count := presets.Get("length").Int()
	for i := 0; i < count; i++ {
		btn := presets.Index(i)
		presetName := btn.Call("getAttribute", "data-room-preset").String()
		btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			setElementText(doc, "currentRoomName", presetName)
			triggerToast("Room Switched", fmt.Sprintf("Joined preset arena: %s", presetName), "info")
			return nil
		}))
	}
}

func openModal(id string) {
	js.Global().Get("SkillNodeBridge").Call("openModal", id)
}

func closeModal(id string) {
	js.Global().Get("SkillNodeBridge").Call("closeModal", id)
}

func triggerToast(title, message, kind string) {
	js.Global().Get("SkillNodeBridge").Call("toast", title, message, kind)
}

func setElementText(doc js.Value, id, text string) {
	el := doc.Call("getElementById", id)
	if !el.IsNull() {
		el.Set("innerText", text)
	}
}

// Temporary stub until Phase 4 typing engine connects
func initTypingView() {}
