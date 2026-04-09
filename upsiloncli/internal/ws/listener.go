package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ecumeurs/upsiloncli/internal/api"
	"github.com/ecumeurs/upsiloncli/internal/display"
	"github.com/ecumeurs/upsiloncli/internal/dto"
	"github.com/ecumeurs/upsiloncli/internal/session"
	"github.com/gorilla/websocket"
)

// Listener manages the real-time WebSocket connection to Laravel Reverb.
type Listener struct {
	Client   *api.Client
	Session  *session.Session
	Printer  *display.Printer
	Conn     *websocket.Conn
	SocketID string
	AppKey   string
	Host     string

	mu   sync.Mutex
	subs map[string]bool
}

// NewListener creates a new WebSocket listener.
func NewListener(client *api.Client, sess *session.Session, printer *display.Printer) *Listener {
	return &Listener{
		Client:  client,
		Session: sess,
		Printer: printer,
		AppKey:  "qtjp54myattne9euwedu", // Hardcoded for this environment
		Host:    "127.0.0.1:8080",      // Hardcoded for this environment
		subs:    make(map[string]bool),
	}
}

// Start opens the connection and starts the message loop.
func (l *Listener) Start() {
	u := fmt.Sprintf("ws://%s/app/%s?protocol=7&client=js&version=8.4.0-rc2&flash=false", l.Host, l.AppKey)

	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		l.Printer.Warn(fmt.Sprintf("WebSocket connection failed: %v", err))
		return
	}
	l.Conn = conn

	go l.listenLoop()
}

func (l *Listener) listenLoop() {
	defer l.Conn.Close()

	for {
		_, message, err := l.Conn.ReadMessage()
		if err != nil {
			return
		}

		var envelope struct {
			Event   string          `json:"event"`
			Channel string          `json:"channel"`
			Data    json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(message, &envelope); err != nil {
			continue
		}

		switch envelope.Event {
		case "pusher:connection_established":
			var data struct {
				SocketID string `json:"socket_id"`
			}
			// Pusher data is sometimes double-encoded in JSON strings
			unquotedData := strings.Trim(string(envelope.Data), "\"")
			if err := json.Unmarshal([]byte(unquotedData), &data); err != nil {
				// Try direct unmarshal if not double-encoded
				json.Unmarshal(envelope.Data, &data)
			}
			l.SocketID = data.SocketID
			l.subscribeToUserChannel()

		case "match.found":
			var data struct {
				MatchID string `json:"match_id"`
			}
			unquoted := strings.Trim(string(envelope.Data), "\"")
			json.Unmarshal([]byte(unquoted), &data)
			
			if data.MatchID != "" {
				l.Session.Set("match_id", data.MatchID)
				l.Printer.WebSocket("MatchFound", envelope.Data)
				l.Printer.System(fmt.Sprintf("Match detected! Initializing arena %s...", data.MatchID))
				
				// Fetch full state (participants + board)
				l.initializeMatch(data.MatchID)
				
				// Subscribe to arena updates
				l.subscribeToArenaChannel(data.MatchID)
			}

		case "board.updated":
			var board dto.BoardState
			unquoted := strings.Trim(string(envelope.Data), "\"")
			if err := json.Unmarshal([]byte(unquoted), &board); err == nil {
				l.Session.SetLastBoard(&board)
				l.Printer.System("Tactical feed updated.")
				// Auto-redraw if board is already displayed? 
				// For now let's just notify. The user can type 'redraw'.
			}

		case "pusher_internal:subscription_succeeded":
			// Handled silently
		}
	}
}

// Sync reconciles active WebSocket subscriptions with the current session state.
// It ensures we are subscribed to the private user channel and any active arena channel.
func (l *Listener) Sync() {
	l.mu.Lock()
	conn := l.Conn
	socketID := l.SocketID
	l.mu.Unlock()

	if conn == nil || socketID == "" {
		return
	}

	// 1. Sync User Channel
	uid := l.Session.UserIdentifier()
	if uid != "" {
		channel := fmt.Sprintf("private-user.%s", uid)
		l.ensureSubscription(channel)
	}

	// 2. Sync Arena Channel
	if mid, ok := l.Session.Get("match_id"); ok && mid != "" {
		channel := fmt.Sprintf("private-arena.%s", mid)
		l.ensureSubscription(channel)
	}
}

func (l *Listener) subscribeToUserChannel() {
	uid := l.Session.UserIdentifier()
	if uid == "" {
		return
	}
	l.ensureSubscription(fmt.Sprintf("private-user.%s", uid))
}

func (l *Listener) subscribeToArenaChannel(matchID string) {
	l.ensureSubscription(fmt.Sprintf("private-arena.%s", matchID))
}

func (l *Listener) ensureSubscription(channel string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Conn == nil || l.SocketID == "" {
		return
	}

	if l.subs[channel] {
		return
	}

	auth, err := l.getAuth(channel)
	if err != nil {
		return
	}

	sub := map[string]interface{}{
		"event": "pusher:subscribe",
		"data": map[string]string{
			"channel": channel,
			"auth":    auth,
		},
	}
	
	if err := l.Conn.WriteJSON(sub); err == nil {
		l.subs[channel] = true
	}
}

func (l *Listener) getAuth(channel string) (string, error) {
	token := l.Session.Token()
	if token == "" {
		return "", fmt.Errorf("not authenticated")
	}

	// POST /api/broadcasting/auth
	url := l.Client.BaseURL + "/api/broadcasting/auth"
	body := strings.NewReader(fmt.Sprintf("socket_id=%s&channel_name=%s", l.SocketID, channel))
	
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Auth string `json:"auth"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Auth, nil
}

func (l *Listener) initializeMatch(matchID string) {
	// Call GET /api/v1/game/{id}
	resp, err := l.Client.Get(fmt.Sprintf("/api/v1/game/%s", matchID))
	if err != nil {
		return
	}

	var game dto.GameResponse
	// We need to re-marshal/unmarshal because resp.Data is interface{}
	dataBytes, _ := json.Marshal(resp.Data)
	if err := json.Unmarshal(dataBytes, &game); err == nil {
		l.Session.SetParticipants(game.Participants)
		l.Session.SetLastBoard(&game.GameState)
	} else {
		// Try unmarshaling from root if it's a direct structured response
		json.Unmarshal([]byte(resp.RawBody), &game)
		l.Session.SetParticipants(game.Participants)
		l.Session.SetLastBoard(&game.GameState)
	}
}
