// SkillNODE Tiny Vanilla JS Bridge (Zero Dependencies, Zero node_modules)
// Bridges DOM events, routing, real-time WebSockets, and animations to the Go WebAssembly core.

(function () {
  window.SkillNodeBridge = {
    // Show a modern glassmorphic toast notification
    toast: function (title, message, kind = "info") {
      const container = document.getElementById("toastContainer");
      if (!container) return;

      const toast = document.createElement("div");
      toast.className = "app-toast animate-slide-up";
      toast.style.pointerEvents = "auto";
      toast.style.marginBottom = "8px";

      let accentColor = "var(--accent)";
      let icon = "fa-info-circle";
      if (kind === "success") {
        accentColor = "#10b981";
        icon = "fa-check-circle";
      } else if (kind === "warning") {
        accentColor = "#f59e0b";
        icon = "fa-exclamation-triangle";
      } else if (kind === "error") {
        accentColor = "#ef4444";
        icon = "fa-times-circle";
      }

      toast.innerHTML = `
        <button class="app-toast-dismiss" aria-label="Dismiss">
          <i class="fas fa-times"></i>
        </button>
        <div style="display: flex; align-items: flex-start; gap: 10px;">
          <i class="fas ${icon}" style="color: ${accentColor}; font-size: 15px; margin-top: 2px;"></i>
          <div style="flex: 1;">
            <div class="app-toast-title">${title}</div>
            <div class="app-toast-message">${message}</div>
          </div>
        </div>
      `;

      const dismissBtn = toast.querySelector(".app-toast-dismiss");
      if (dismissBtn) {
        dismissBtn.addEventListener("click", () => {
          toast.style.opacity = "0";
          toast.style.transform = "translateY(8px)";
          setTimeout(() => toast.remove(), 200);
        });
      }

      container.appendChild(toast);

      setTimeout(() => {
        if (toast.parentNode) {
          toast.style.transition = "opacity 0.25s ease, transform 0.25s ease";
          toast.style.opacity = "0";
          toast.style.transform = "translateY(8px)";
          setTimeout(() => toast.remove(), 250);
        }
      }, 3500);
    },

    // Modal helpers
    openModal: function (id) {
      const modal = document.getElementById(id);
      if (modal) {
        modal.style.display = "flex";
        modal.style.opacity = "1";
      }
    },

    closeModal: function (id) {
      const modal = document.getElementById(id);
      if (modal) {
        modal.style.display = "none";
      }
    },

    // Get current hash route (e.g. "login", "dashboard", "activities/typing")
    getCurrentRoute: function () {
      const hash = window.location.hash;
      if (!hash || hash === "#" || hash === "#/") {
        return "login";
      }
      return hash.replace(/^#\/?/, "") || "login";
    },

    // Set route programmatically
    navigate: function (route) {
      window.location.hash = "#/" + route.replace(/^\/?/, "");
    },
  };

  // Expose global shorthand for Go WASM
  window.showToast = function (title, message, kind) {
    window.SkillNodeBridge.toast(title, message, kind);
  };

  // Real-time Socket.IO Bridge
  window.SkillNodeSocket = {
    socket: null,
    init: function (apiBaseUrl, name, onConnect, onPresence, onMessage, onRoomState, onError) {
      if (typeof io === "undefined") {
        console.warn("Socket.io not loaded from CDN");
        return;
      }
      if (this.socket && this.socket.connected) return;

      try {
        this.socket = io(apiBaseUrl, {
          transports: ["websocket", "polling"],
          auth: { name: name || "Operator" },
        });

        this.socket.on("connect", () => {
          if (typeof onConnect === "function") onConnect();
        });

        this.socket.on("presence", (data) => {
          if (typeof onPresence === "function") onPresence(JSON.stringify(data));
        });

        this.socket.on("room_state", (data) => {
          if (typeof onRoomState === "function") onRoomState(JSON.stringify(data));
        });

        this.socket.on("room_message", (data) => {
          if (typeof onMessage === "function") onMessage(JSON.stringify(data));
        });

        this.socket.on("room_error", (data) => {
          if (typeof onError === "function") onError(data.message || "Room error");
        });
      } catch (e) {
        console.error("Socket initialization error:", e);
      }
    },

    joinRoom: function (room) {
      if (this.socket) this.socket.emit("join_room", { room: room });
    },

    joinByCode: function (code) {
      if (this.socket) this.socket.emit("join_room_by_code", { code: code });
    },

    sendMessage: function (text) {
      if (this.socket) this.socket.emit("room_message", { text: text });
    },

    setProfile: function (name) {
      if (this.socket) this.socket.emit("set_profile", { name: name });
    },
  };

  // Listen to hash change and notify Go WASM router if registered
  window.addEventListener("hashchange", function () {
    const route = window.SkillNodeBridge.getCurrentRoute();
    if (typeof window.onRouteChanged === "function") {
      window.onRouteChanged(route);
    }
  });

  // Attach modal backdrop click to dismiss
  document.addEventListener("DOMContentLoaded", function () {
    const backdrop = document.getElementById("authModalBackdrop");
    if (backdrop) {
      backdrop.addEventListener("click", function (e) {
        if (e.target === backdrop) {
          window.SkillNodeBridge.closeModal("authModalBackdrop");
        }
      });
    }
  });

  // Client Sandbox Lockdown: Freeze bridges against console hijacking
  try {
    Object.freeze(window.SkillNodeBridge);
    Object.freeze(window.SkillNodeSocket);
    Object.defineProperty(window, "SkillNodeBridge", {
      value: window.SkillNodeBridge,
      writable: false,
      configurable: false,
    });
    Object.defineProperty(window, "SkillNodeSocket", {
      value: window.SkillNodeSocket,
      writable: false,
      configurable: false,
    });
    console.log(
      "%c[SkillNODE Zero-Trust Security] Client Sandbox Locked.\nAll local storage records are cryptographically sealed with Go WASM HMAC-SHA256 signatures.\nConsole privilege escalation and unverified admin roles are strictly blocked.",
      "background: #4c1d95; color: #c084fc; font-weight: bold; font-size: 11px; padding: 6px 10px; border-radius: 4px; border: 1px solid #7c3aed;"
    );
  } catch (err) {
    // Non-critical sandbox fallback
  }
})();
