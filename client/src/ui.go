//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

var currentActiveRoute = "login"

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
	if cleanRoute == "" || !IsSessionAuthenticated() {
		cleanRoute = "login"
	}
	SwitchView(cleanRoute)

	// 5. Theme toggle button in sidebar
	themeBtn := doc.Call("getElementById", "toggleThemeBtn")
	if !themeBtn.IsNull() {
		themeBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			toggleTheme()
			return nil
		}))
	}

	// 6. Logout trigger in sidebar
	logoutBtn := doc.Call("getElementById", "logoutTriggerBtn")
	if !logoutBtn.IsNull() {
		logoutBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			LockSession()
			renderProfileInfo()
			triggerToast("Session Locked", "Signed out of SkillNODE arena.", "info")
			js.Global().Get("location").Set("hash", "#/login")
			SwitchView("login")
			return nil
		}))
	}

	// 7. Modal close button
	authCloseBtn := doc.Call("getElementById", "authModalCloseBtn")
	if !authCloseBtn.IsNull() {
		authCloseBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			closeModal("authModalBackdrop")
			return nil
		}))
	}

	// 8. Profile view edit modal trigger & share
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

	// 9. Bind room presets in multiplayer view
	bindRoomPresetButtons()

	// 10. Bind landing page events (tabs, forms, demo, smooth scrolling)
	bindLandingEvents()
}

func bindLandingEvents() {
	doc := js.Global().Get("document")

	// Smooth scroll helper
	scrollFn := func(targetID string) js.Func {
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 && !args[0].IsNull() && !args[0].IsUndefined() {
				args[0].Call("preventDefault")
			}
			scrollToElement(targetID)
			return nil
		})
	}

	// Navigation & Banner scroll links
	bindClick(doc, "landingBannerTrackLink", scrollFn("tracks"))
	bindClick(doc, "navLinkTracks", scrollFn("tracks"))
	bindClick(doc, "heroTracksBtn", scrollFn("tracks"))
	bindClick(doc, "navLinkFeatures", scrollFn("features"))
	bindClick(doc, "navLinkBenchmarks", scrollFn("benchmarks"))
	bindClick(doc, "landingNavGetStartedBtn", scrollFn("auth-portal"))
	bindClick(doc, "heroEnterArenaBtn", scrollFn("auth-portal"))
	bindClick(doc, "bottomLaunchBtn", scrollFn("auth-portal"))

	// Landing theme toggle button
	landingThemeBtn := doc.Call("getElementById", "landingThemeBtn")
	if !landingThemeBtn.IsNull() {
		landingThemeBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			toggleTheme()
			return nil
		}))
	}

	// Tab switching: Sign In vs Create Account
	tabSignIn := doc.Call("getElementById", "landingTabSignIn")
	if !tabSignIn.IsNull() {
		tabSignIn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			setLandingAuthMode("signin")
			return nil
		}))
	}

	tabSignUp := doc.Call("getElementById", "landingTabSignUp")
	if !tabSignUp.IsNull() {
		tabSignUp.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			setLandingAuthMode("signup")
			return nil
		}))
	}

	// Password visibility toggles
	bindPasswordToggle("landingLoginTogglePwd", "landingLoginPassword", "landingLoginPwdIcon")
	bindPasswordToggle("landingSignupTogglePwd", "landingSignupPassword", "landingSignupPwdIcon")

	// Instant Demo Access (0-Click Entry)
	demoBtn := doc.Call("getElementById", "landingDemoBtn")
	if !demoBtn.IsNull() {
		demoBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			LoginWithDemo()
			renderProfileInfo()
			triggerToast("Demo Access Granted", "Entering SkillNODE as operator Alex Coder.", "success")
			js.Global().Get("location").Set("hash", "#/dashboard")
			SwitchView("dashboard")
			return nil
		}))
	}

	// Sign In Form Submission
	signInForm := doc.Call("getElementById", "landingSignInForm")
	if !signInForm.IsNull() {
		signInForm.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				args[0].Call("preventDefault")
			}
			uVal := doc.Call("getElementById", "landingLoginUsername").Get("value").String()
			pVal := doc.Call("getElementById", "landingLoginPassword").Get("value").String()

			_, err := LoginWithCredentials(uVal, pVal)
			if err != nil {
				triggerToast("Authentication Failed", err.Error(), "error")
				return nil
			}

			renderProfileInfo()
			triggerToast("Welcome back!", "Authenticated into SkillNODE arena.", "success")
			js.Global().Get("location").Set("hash", "#/dashboard")
			SwitchView("dashboard")
			return nil
		}))
	}

	// Sign Up Form Submission
	signUpForm := doc.Call("getElementById", "landingSignUpForm")
	if !signUpForm.IsNull() {
		signUpForm.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				args[0].Call("preventDefault")
			}
			name := doc.Call("getElementById", "landingSignupName").Get("value").String()
			username := doc.Call("getElementById", "landingSignupUsername").Get("value").String()
			email := doc.Call("getElementById", "landingSignupEmail").Get("value").String()
			region := doc.Call("getElementById", "landingSignupRegion").Get("value").String()

			_, err := CreateAccountProfile(name, username, email, region)
			if err != nil {
				triggerToast("Sign-up Failed", err.Error(), "error")
				return nil
			}

			renderProfileInfo()
			triggerToast("Welcome to SkillNODE", "Your professional skill profile is ready.", "success")
			js.Global().Get("location").Set("hash", "#/dashboard")
			SwitchView("dashboard")
			return nil
		}))
	}
}

func bindClick(doc js.Value, id string, fn js.Func) {
	el := doc.Call("getElementById", id)
	if !el.IsNull() {
		el.Call("addEventListener", "click", fn)
	}
}

func scrollToElement(id string) {
	doc := js.Global().Get("document")
	el := doc.Call("getElementById", id)
	if !el.IsNull() {
		opt := js.Global().Get("Object").New()
		opt.Set("behavior", "smooth")
		opt.Set("block", "start")
		el.Call("scrollIntoView", opt)
	}
}

func setLandingAuthMode(mode string) {
	doc := js.Global().Get("document")
	tabSignIn := doc.Call("getElementById", "landingTabSignIn")
	tabSignUp := doc.Call("getElementById", "landingTabSignUp")
	formSignIn := doc.Call("getElementById", "landingSignInForm")
	formSignUp := doc.Call("getElementById", "landingSignUpForm")
	title := doc.Call("getElementById", "landingAuthTitle")

	if mode == "signin" {
		if !tabSignIn.IsNull() {
			tabSignIn.Get("style").Set("background", "var(--btn-bg)")
			tabSignIn.Get("style").Set("color", "#fff")
		}
		if !tabSignUp.IsNull() {
			tabSignUp.Get("style").Set("background", "transparent")
			tabSignUp.Get("style").Set("color", "var(--text-muted)")
		}
		if !formSignIn.IsNull() {
			formSignIn.Get("style").Set("display", "flex")
		}
		if !formSignUp.IsNull() {
			formSignUp.Get("style").Set("display", "none")
		}
		if !title.IsNull() {
			title.Set("innerText", "Sign in to SkillNODE")
		}
	} else {
		if !tabSignIn.IsNull() {
			tabSignIn.Get("style").Set("background", "transparent")
			tabSignIn.Get("style").Set("color", "var(--text-muted)")
		}
		if !tabSignUp.IsNull() {
			tabSignUp.Get("style").Set("background", "var(--btn-bg)")
			tabSignUp.Get("style").Set("color", "#fff")
		}
		if !formSignIn.IsNull() {
			formSignIn.Get("style").Set("display", "none")
		}
		if !formSignUp.IsNull() {
			formSignUp.Get("style").Set("display", "flex")
		}
		if !title.IsNull() {
			title.Set("innerText", "Create Operator Profile")
		}
	}
}

func bindPasswordToggle(btnID, inputID, iconID string) {
	doc := js.Global().Get("document")
	btn := doc.Call("getElementById", btnID)
	input := doc.Call("getElementById", inputID)
	icon := doc.Call("getElementById", iconID)
	if btn.IsNull() || input.IsNull() {
		return
	}
	btn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		curr := input.Get("type").String()
		if curr == "password" {
			input.Set("type", "text")
			if !icon.IsNull() {
				icon.Set("className", "fas fa-eye-slash")
			}
		} else {
			input.Set("type", "password")
			if !icon.IsNull() {
				icon.Set("className", "fas fa-eye")
			}
		}
		return nil
	}))
}

func SwitchView(route string) {
	route = strings.TrimPrefix(route, "#/")
	route = strings.TrimPrefix(route, "#")
	route = strings.TrimPrefix(route, "/")

	doc := js.Global().Get("document")
	landingView := doc.Call("getElementById", "view-landing")
	appShell := doc.Call("getElementById", "root")

	// If visitor is unauthenticated or explicitly requesting the login / landing route
	if route == "login" || route == "" || !IsSessionAuthenticated() {
		currentActiveRoute = "login"
		if !landingView.IsNull() {
			landingView.Get("style").Set("display", "flex")
			landingView.Set("scrollTop", 0)
		}
		if !appShell.IsNull() {
			appShell.Get("style").Set("display", "none")
		}
		doc.Set("title", "SkillNODE")
		return
	}

	// User is authenticated and navigating inside the arena
	currentActiveRoute = route
	if !landingView.IsNull() {
		landingView.Get("style").Set("display", "none")
	}
	if !appShell.IsNull() {
		appShell.Get("style").Set("display", "flex")
	}
	doc.Set("title", "SkillNODE")

	// Hide all views inside pageFrame
	views := doc.Call("querySelectorAll", "#pageFrame .page-view")
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
		initCodingView()
	case strings.HasPrefix(route, "activities/math"):
		targetViewID = "view-activity-math"
		viewTitle = "Quant Lab: Math Clash"
		navKey = "activities"
		initMathView()
	case strings.HasPrefix(route, "multiplayer"):
		targetViewID = "view-multiplayer"
		viewTitle = "Battle Net: Real-Time Multiplayer"
		navKey = "multiplayer"
		initMultiplayerView()
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

	landingIcon := doc.Call("getElementById", "landingThemeIcon")
	landingLabel := doc.Call("getElementById", "landingThemeLabel")
	if !landingIcon.IsNull() {
		if strings.Contains(theme, "light") {
			landingIcon.Set("className", "fas fa-sun")
		} else {
			landingIcon.Set("className", "fas fa-moon")
		}
	}
	if !landingLabel.IsNull() {
		if strings.Contains(theme, "light") {
			landingLabel.Set("innerText", "Light")
		} else {
			landingLabel.Set("innerText", "Dark")
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
	applyStoredTheme()

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
	setElementText(doc, "profileRoleDisplay", strings.ToUpper(p.Role)+" (LOCAL SANDBOX)")

	badgeEl := doc.Call("getElementById", "profileBadgeDisplay")
	if !badgeEl.IsNull() {
		if p.IsVerified {
			badgeEl.Set("innerHTML", `<i class="fas fa-check-circle"></i> VERIFIED OPERATOR`)
			badgeEl.Get("style").Set("background", "rgba(16, 185, 129, 0.15)")
			badgeEl.Get("style").Set("color", "#34d399")
		} else {
			badgeEl.Set("innerHTML", `<i class="fas fa-shield-alt"></i> GUEST OPERATOR`)
			badgeEl.Get("style").Set("background", "rgba(147, 51, 234, 0.15)")
			badgeEl.Get("style").Set("color", "var(--accent-bright)")
		}
	}

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
