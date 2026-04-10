// @spec-link [[api_websocket_game_events]]
package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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

	waitMu  sync.Mutex
	waiters map[string][]chan interface{}
}

// NewListener creates a new WebSocket listener.
func NewListener(client *api.Client, sess *session.Session, printer *display.Printer) *Listener {
	l := &Listener{
		Client:  client,
		Session: sess,
		Printer: printer,
		AppKey:  os.Getenv("REVERB_APP_KEY"),
		Host:    os.Getenv("REVERB_HOST"),
		subs:    make(map[string]bool),
		waiters: make(map[string][]chan interface{}),
	}
	if l.Host == "" {
		l.Host = "127.0.0.1:8080"
	}
	return l
}

// Start opens the connection and starts the message loop.
func (l *Listener) Start() {
	u := fmt.Sprintf("ws://%s/app/%s?protocol=7&client=js&version=8.4.0-rc2&flash=false", l.Host, l.AppKey)
	if l.Printer != nil {
		l.Printer.Wscat(u)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		if l.Printer != nil {
			l.Printer.Warn(fmt.Sprintf("WebSocket connection failed (is Reverb running?): %v", err))
		}
		return
	}
	l.Conn = conn
	if l.Printer != nil {
		l.Printer.System("WebSocket link established. Waiting for handshake...")
	}

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
			// Pusher data is sometimes double-encoded as a JSON string
			var dataStr string
			if err := json.Unmarshal(envelope.Data, &dataStr); err == nil {
				json.Unmarshal([]byte(dataStr), &data)
			} else {
				json.Unmarshal(envelope.Data, &data)
			}
			
			l.SocketID = data.SocketID
			if l.Printer != nil {
				l.Printer.System(fmt.Sprintf("Handshake successful. SocketID: %s", l.SocketID))
			}
			l.subscribeToUserChannel()
			l.Sync()

		case "match.found":
			// Laravel payload: { match_id: "...", user_id: "...", data: [] }
			var payload struct {
				MatchID string `json:"match_id"`
			}
			
			// Reverb/Pusher data is sometimes double-encoded as a JSON string
			var dataStr string
			if err := json.Unmarshal(envelope.Data, &dataStr); err == nil {
				json.Unmarshal([]byte(dataStr), &payload)
			} else {
				json.Unmarshal(envelope.Data, &payload)
			}
			
			if payload.MatchID != "" {
				l.Session.Set("match_id", payload.MatchID)
				if l.Printer != nil {
					l.Printer.WebSocket("MatchFound", envelope.Data)
					l.Printer.System(fmt.Sprintf("Match detected! Initializing arena %s...", payload.MatchID))
				}
				
				// Fetch full state (participants + board)
				l.initializeMatch(payload.MatchID)
				
				// Subscribe to arena updates
				l.subscribeToArenaChannel(payload.MatchID)
			} else {
				if l.Printer != nil {
					l.Printer.Warn(fmt.Sprintf("Received match.found but match_id is empty. Raw: %s", string(envelope.Data)))
				}
			}

		case "board.updated":
			// Laravel payload: { match_id: "...", data: { ...board... } }
			var payload struct {
				Data dto.BoardState `json:"data"`
			}

			if l.Printer != nil {
				l.Printer.WebSocket("board.updated", envelope.Data)
			}

			var dataStr string
			if err := json.Unmarshal(envelope.Data, &dataStr); err == nil {
				if err := json.Unmarshal([]byte(dataStr), &payload); err == nil {
					l.Session.SetLastBoard(&payload.Data)
					if l.Printer != nil {
						l.Printer.System("Tactical feed updated.")
						l.Printer.Suggestions([]string{"redraw"})
					}
				} else {
					if l.Printer != nil {
						l.Printer.Warn(fmt.Sprintf("Failed to decode board.updated data string: %v", err))
					}
				}
			} else {
				if err := json.Unmarshal(envelope.Data, &payload); err == nil {
					l.Session.SetLastBoard(&payload.Data)
					if l.Printer != nil {
						l.Printer.System("Tactical feed updated.")
						l.Printer.Suggestions([]string{"redraw"})
					}
				} else {
					if l.Printer != nil {
						l.Printer.Warn(fmt.Sprintf("Failed to decode board.updated payload: %v", err))
					}
				}
			}

		case "pusher_internal:subscription_succeeded":
			if l.Printer != nil {
				l.Printer.System(fmt.Sprintf("Subscription for %s acknowledged by server.", envelope.Channel))
			}
		
		case "pusher:ping":
			// Respond to server heartbeats to prevent timeout (Error 4201)
			l.Conn.WriteJSON(map[string]string{"event": "pusher:pong"})
		
		default:
			// Print all other events for transparency as requested
			if l.Printer != nil {
				l.Printer.WebSocket(envelope.Event, envelope.Data)
			}
		}

		// Notify any waiters for this event
		l.notifyWaiters(envelope.Event, envelope.Data)
	}
}

// WaitForData blocks until an event of the given name is received or timeout occurs.
func (l *Listener) WaitForData(eventName string, timeoutMs int) (interface{}, error) {
	ch := make(chan interface{}, 1)
	
	l.waitMu.Lock()
	l.waiters[eventName] = append(l.waiters[eventName], ch)
	l.waitMu.Unlock()

	defer func() {
		l.waitMu.Lock()
		defer l.waitMu.Unlock()
		// Remove ch from waiters
		list := l.waiters[eventName]
		for i, v := range list {
			if v == ch {
				l.waiters[eventName] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}()

	select {
	case data := <-ch:
		return data, nil
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		return nil, fmt.Errorf("timeout waiting for event: %s", eventName)
	}
}

func (l *Listener) notifyWaiters(eventName string, data json.RawMessage) {
	l.waitMu.Lock()
	defer l.waitMu.Unlock()

	waiters, ok := l.waiters[eventName]
	if !ok || len(waiters) == 0 {
		return
	}

	// Parse data into interface{} so it's clean for JS
	var parsed interface{}
	// Reverb/Pusher data is sometimes double-encoded as a JSON string
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err == nil {
		json.Unmarshal([]byte(dataStr), &parsed)
	} else {
		json.Unmarshal(data, &parsed)
	}

	for _, ch := range waiters {
		select {
		case ch <- parsed:
		default: // skip if channel is full
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

// Status returns the current health of the listener.
func (l *Listener) Status() (connected bool, socketID string, subscriptions []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	connected = l.Conn != nil
	socketID = l.SocketID
	subscriptions = make([]string, 0, len(l.subs))
	for sub := range l.subs {
		subscriptions = append(subscriptions, sub)
	}
	return
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

	// broadcasting/auth (no /api prefix)
	url := l.Client.BaseURL + "/broadcasting/auth"
	formBody := fmt.Sprintf("socket_id=%s&channel_name=%s", l.SocketID, channel)
	
	// Display the manual test command as requested
	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Accept", "application/json") // Explicitly request JSON
	headers.Set("Authorization", "Bearer "+token)
	if l.Printer != nil {
		l.Printer.Curl("POST", url, headers, []byte(formBody))
	}

	body := strings.NewReader(formBody)
	req, _ := http.NewRequest("POST", url, body)
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if l.Printer != nil {
			l.Printer.Warn(fmt.Sprintf("Authorization failed for %s (Status %d)", channel, resp.StatusCode))
			l.Printer.Warn(fmt.Sprintf("Raw Body: %s", string(raw)))
		}
		return "", fmt.Errorf("auth failed: %d", resp.StatusCode)
	}

	var result struct {
		Auth string `json:"auth"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		if l.Printer != nil {
			l.Printer.Warn(fmt.Sprintf("Failed to decode auth response: %v", err))
		}
		return "", err
	}

	// Display the manual wscat payload as requested
	if l.Printer != nil {
		l.Printer.WscatPayload(channel, result.Auth)
	}

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
