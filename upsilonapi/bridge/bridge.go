package bridge

// @spec-link [[module_upsilonapi]]

import (
	"fmt"
	"sync"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/ruler"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/google/uuid"
)

type ArenaBridge struct {
	mu     sync.RWMutex
	arenas map[uuid.UUID]*ruler.Ruler
}

var bridge = &ArenaBridge{
	arenas: make(map[uuid.UUID]*ruler.Ruler),
}

func Get() *ArenaBridge {
	return bridge
}

func (b *ArenaBridge) StartArena(callbackURL string) (uuid.UUID, *rulermethods.AddControllerReply, error) {
	r := ruler.NewRuler()
	b.mu.Lock()
	b.arenas[r.ID] = &r
	b.mu.Unlock()

	// We need at least one controller to get the initial state
	// In the future, we might add multiple based on players payload
	hc := NewHTTPController(callbackURL)

	msg := message.Create(hc, rulermethods.AddController{
		Controller:   hc,
		ControllerID: hc.ID,
	}, nil)

	// We need to wait for the reply to get the initial state
	respChan := make(chan *message.Message, 1)
	r.SendActor(msg, respChan)

	// Wait for response or timeout
	select {
	case m := <-respChan:
		if m.HasError {
			return uuid.Nil, nil, fmt.Errorf("failed to add controller: %s", m.ErrorMessage)
		}
		reply := m.Content.(rulermethods.AddControllerReply)
		return r.ID, &reply, nil
	case <-time.After(5 * time.Second):
		return uuid.Nil, nil, fmt.Errorf("timeout waiting for ruler response")
	}
}

func (b *ArenaBridge) GetArena(id uuid.UUID) (*ruler.Ruler, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	r, ok := b.arenas[id]
	return r, ok
}
