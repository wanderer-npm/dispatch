package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/wanderer-npm/dispatch/internal/config"
	"github.com/wanderer-npm/dispatch/internal/discord"
	"github.com/wanderer-npm/dispatch/internal/events"
)

type eventFn func([]byte) (*discord.Embed, error)

var registry = map[string]eventFn{
	"push":         events.Push,
	"repository":   events.Repository,
	"create":       events.Create,
	"delete":       events.Delete,
	"fork":         events.Fork,
	"watch":        events.Star,
	"pull_request": events.PullRequest,
	"release":      events.Release,
	"member":       events.Member,
	"issues":       events.Issues,
}

type webhookHandler struct {
	cfg *config.Config
}

func (h *webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	if !verifySignature(body, r.Header.Get("X-Hub-Signature-256"), h.cfg.Server.Secret) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "ping" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	fn, ok := registry[eventType]
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	embed, err := fn(body)
	if err != nil {
		log.Printf("[dispatch] event=%s err=%v", eventType, err)
		http.Error(w, "failed to process event", http.StatusInternalServerError)
		return
	}
	if embed == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for _, url := range h.cfg.WebhooksFor(eventType) {
		if err := discord.Send(url, *embed); err != nil {
			log.Printf("[dispatch] discord send failed url=%s err=%v", url, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func verifySignature(body []byte, signature, secret string) bool {
	if secret == "" {
		return true
	}
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
