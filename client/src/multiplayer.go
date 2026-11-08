//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
)

type SocketMember struct {
	SocketID string `json:"socketId"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
}

type SocketMessage struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type SocketPresence struct {
	Room    string         `json:"room"`
	Members []SocketMember `json:"members"`
}

type SocketRoomState struct {
	Room     string          `json:"room"`
	JoinCode string          `json:"joinCode"`
	Members  []SocketMember  `json:"members"`
	Messages []SocketMessage `json:"messages"`
}

var multiplayerInitialized = false
var currentRoomCode = ""

func initMultiplayerView() {
	if multiplayerInitialized {
		return
	}
	multiplayerInitialized = true

	doc := js.Global().Get("document")
	p := GetSessionProfile()
	playerName := "Operator"
	if p != nil && p.Name != "" {
		playerName = p.Name
	}

	apiURL := getAPIBaseURL()

	// 1. Initialize Socket.IO bridge with JS callbacks
	socketBridge := js.Global().Get("SkillNodeSocket")
	if !socketBridge.IsUndefined() {
		onConnect := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			statusEl := doc.Call("getElementById", "multiplayerSocketStatus")
			if !statusEl.IsNull() {
				statusEl.Set("innerText", "Realtime Synchronized (Connected)")
			}
			socketBridge.Call("joinRoom", "lobby")
			return nil
		})

		onPresence := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				var pres SocketPresence
				if err := json.Unmarshal([]byte(args[0].String()), &pres); err == nil {
					renderPeerMembers(pres.Members)
				}
			}
			return nil
		})

		onRoomState := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				var state SocketRoomState
				if err := json.Unmarshal([]byte(args[0].String()), &state); err == nil {
					currentRoomCode = state.JoinCode
					setElementText(doc, "currentRoomName", state.Room)
					renderPeerMembers(state.Members)
					renderRoomChatHistory(state.Messages)
				}
			}
			return nil
		})

		onMessage := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				var msg SocketMessage
				if err := json.Unmarshal([]byte(args[0].String()), &msg); err == nil {
					appendSingleChatMessage(msg)
				}
			}
			return nil
		})

		onError := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				triggerToast("Room Notification", args[0].String(), "warning")
			}
			return nil
		})

		socketBridge.Call("init", apiURL, playerName, onConnect, onPresence, onMessage, onRoomState, onError)
	}

	// 2. Chat Send Button listener
	sendBtn := doc.Call("getElementById", "multiplayerSendChatBtn")
	if !sendBtn.IsNull() {
		sendBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			sendChatFromInput()
			return nil
		}))
	}

	// 3. Chat Input Enter key listener
	chatInput := doc.Call("getElementById", "multiplayerChatInput")
	if !chatInput.IsNull() {
		chatInput.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 && args[0].Get("key").String() == "Enter" {
				sendChatFromInput()
			}
			return nil
		}))
	}

	// 4. Join by code button
	joinCodeBtn := doc.Call("getElementById", "multiplayerJoinCodeBtn")
	if !joinCodeBtn.IsNull() {
		joinCodeBtn.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			codeInput := doc.Call("getElementById", "multiplayerCodeInput")
			if !codeInput.IsNull() {
				code := strings.TrimSpace(codeInput.Get("value").String())
				if code != "" {
					js.Global().Get("SkillNodeSocket").Call("joinByCode", code)
					codeInput.Set("value", "")
					triggerToast("Joining Room", fmt.Sprintf("Connecting with code %s...", code), "info")
				}
			}
			return nil
		}))
	}
}

func sendChatFromInput() {
	doc := js.Global().Get("document")
	chatInput := doc.Call("getElementById", "multiplayerChatInput")
	if chatInput.IsNull() {
		return
	}

	text := strings.TrimSpace(chatInput.Get("value").String())
	if text == "" {
		return
	}

	js.Global().Get("SkillNodeSocket").Call("sendMessage", text)
	chatInput.Set("value", "")
}

func renderPeerMembers(members []SocketMember) {
	doc := js.Global().Get("document")

	// Update count badge
	countBadge := doc.Call("getElementById", "multiplayerMemberCountBadge")
	if !countBadge.IsNull() {
		countBadge.Set("innerText", fmt.Sprintf("%d Online", len(members)))
	}

	// Update list
	listEl := doc.Call("getElementById", "multiplayerMemberList")
	if listEl.IsNull() {
		return
	}

	var sb strings.Builder
	p := GetSessionProfile()
	myID := ""
	if p != nil {
		myID = p.ID
	}

	for _, m := range members {
		isMe := m.UserID == myID || m.Name == p.Name
		letter := "O"
		if len(m.Name) > 0 {
			letter = strings.ToUpper(string(m.Name[0]))
		}

		suffix := ""
		if isMe {
			suffix = " (You)"
		}

		sb.WriteString(fmt.Sprintf(`
			<div class="multiplayer-member-row">
				<div style="display: flex; align-items: center; gap: 8px;">
					<div style="width: 28px; height: 28px; border-radius: 50%; background: var(--accent); color: #fff; font-size: 0.75rem; font-weight: 700; display: flex; align-items: center; justify-content: center;">
						%s
					</div>
					<span style="font-size: 0.85rem; font-weight: 600; color: var(--text);">%s%s</span>
				</div>
				<span class="multiplayer-member-chip">Active</span>
			</div>
		`, letter, m.Name, suffix))
	}

	listEl.Set("innerHTML", sb.String())
}

func renderRoomChatHistory(messages []SocketMessage) {
	doc := js.Global().Get("document")
	feed := doc.Call("getElementById", "multiplayerChatFeed")
	if feed.IsNull() {
		return
	}

	p := GetSessionProfile()
	myID := ""
	if p != nil {
		myID = p.ID
	}

	var sb strings.Builder
	for _, m := range messages {
		isMine := m.UserID == myID || m.Author == p.Name
		mineClass := ""
		if isMine {
			mineClass = " mine"
		}

		sb.WriteString(fmt.Sprintf(`
			<div class="multiplayer-message%s">
				<div style="font-size: 0.72rem; color: var(--text-faint); margin-bottom: 2px;">%s</div>
				<div style="font-size: 0.85rem; color: var(--text);">%s</div>
			</div>
		`, mineClass, m.Author, m.Text))
	}

	feed.Set("innerHTML", sb.String())
	feed.Set("scrollTop", feed.Get("scrollHeight"))
}

func appendSingleChatMessage(m SocketMessage) {
	doc := js.Global().Get("document")
	feed := doc.Call("getElementById", "multiplayerChatFeed")
	if feed.IsNull() {
		return
	}

	p := GetSessionProfile()
	myID := ""
	if p != nil {
		myID = p.ID
	}

	isMine := m.UserID == myID || m.Author == p.Name
	mineClass := ""
	if isMine {
		mineClass = " mine"
	}

	msgDiv := doc.Call("createElement", "div")
	msgDiv.Set("className", fmt.Sprintf("multiplayer-message%s animate-slide-up", mineClass))
	msgDiv.Set("innerHTML", fmt.Sprintf(`
		<div style="font-size: 0.72rem; color: var(--text-faint); margin-bottom: 2px;">%s</div>
		<div style="font-size: 0.85rem; color: var(--text);">%s</div>
	`, m.Author, m.Text))

	feed.Call("appendChild", msgDiv)
	feed.Set("scrollTop", feed.Get("scrollHeight"))
}

func getAPIBaseURL() string {
	origin := js.Global().Get("location").Get("origin").String()
	if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
		return "http://localhost:5000"
	}
	return "https://skillnode-api.onrender.com"
}
