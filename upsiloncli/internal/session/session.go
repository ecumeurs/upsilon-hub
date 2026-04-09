// Package session manages JWT tokens and contextual state
// accumulated from API responses during a CLI session.
package session

import (
	"fmt"
	"sync"
)

// Session holds the active JWT and a key-value context store
// populated from API response data (user_id, match_id, etc.).
type Session struct {
	mu      sync.RWMutex
	token   string
	context map[string]string
}

// New creates an empty session.
func New() *Session {
	return &Session{
		context: make(map[string]string),
	}
}

// --- JWT Management ---

// SetToken stores a new JWT. Called after login/register or renewal.
func (s *Session) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
}

// Token returns the current JWT (empty string if unauthenticated).
func (s *Session) Token() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

// ClearToken wipes the JWT. Called on logout or account deletion.
func (s *Session) ClearToken() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
}

// HasToken returns true if a JWT is currently cached.
func (s *Session) HasToken() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token != ""
}

// --- Context Store ---

// Set stores a named value in the session context.
// Values are typically extracted from API responses (e.g., "user_id", "match_id").
func (s *Session) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.context[key] = value
}

// Get retrieves a value from the session context.
// Returns the value and whether it was found.
func (s *Session) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.context[key]
	return v, ok
}

// Delete removes a key from the session context.
func (s *Session) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.context, key)
}

// Clear wipes the entire context (but preserves the JWT).
func (s *Session) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.context = make(map[string]string)
}

// ClearAll wipes both the JWT and the context.
func (s *Session) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
	s.context = make(map[string]string)
}

// Dump returns a snapshot of the session for display.
func (s *Session) Dump() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.context)+1)
	if s.token != "" {
		// Show only the last 12 chars for security
		if len(s.token) > 12 {
			out["jwt"] = "..." + s.token[len(s.token)-12:]
		} else {
			out["jwt"] = s.token
		}
	} else {
		out["jwt"] = "(none)"
	}
	for k, v := range s.context {
		out[k] = v
	}
	return out
}

// HandleTokenRenewal checks a response envelope for meta.token
// and transparently rotates the JWT if present.
// Returns true if a renewal occurred.
func (s *Session) HandleTokenRenewal(meta map[string]interface{}) bool {
	if meta == nil {
		return false
	}
	tokenRaw, ok := meta["token"]
	if !ok {
		return false
	}
	newToken, ok := tokenRaw.(string)
	if !ok || newToken == "" {
		return false
	}
	s.SetToken(newToken)
	return true
}

// String returns a concise session summary for the prompt.
func (s *Session) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, _ := s.context["user_id"]
	matchID, _ := s.context["match_id"]

	auth := "✗"
	if s.token != "" {
		auth = "✓"
	}

	return fmt.Sprintf("auth:%s user:%s match:%s", auth, valueOrDash(userID), valueOrDash(matchID))
}

func valueOrDash(v string) string {
	if v == "" {
		return "-"
	}
	if len(v) > 8 {
		return v[:8]
	}
	return v
}
