# UI Checks


fonts.googleapis.com…;600&display=swap:1 
 Failed to load resource: net::ERR_NAME_NOT_RESOLVED

EngagementHub.vue:95 Match Found via WebSocket! 
Object


BattleArena.vue:88 [Arena] Board empty, attempting fallback sync...
---
waited 1 minute before board got accessible on front. inacceptable: investigate

---
BattleArena.vue:45 [BoardUpdated] Received payload: 
{match_id: '019da960-d66a-7034-848f-a8b51b8caedb', players: Array(2), grid: {…}, turn: Array(5), current_entity_id: 'c7288218-b119-4e41-9ada-fa2db54c9dd9', …}


---

observed: two entity spawnned at the same tile.

----

is shotclock timeout date provided by the go api ? (how does the Front UI now when the shotclock will be triggered ?)

---

front should be notified when go api crash (because it kills the game).
Might need a heartbeat toward the go api. (update communication.md and relevant atd as well)
Create a new issue: attempting game resurrection from board state.

---

add a go api endpoint: does match exists ? if the endpoint doesn't exist already. 
when login in a player is auto forwarded toward a match, should the match not exists, he should be errored(some session crashed stuff) back toward dashboard, and the match data removed from the database (it's faulty anyway)
CLI game_state should return with error session crashed. 

---

investigate logs stored at upsilonapi_crash: upsilon api did crash during the game. 

---

for some reason leaderboard is still empty.

---

front reroll trigger a js alert instead of a html modal, but works.

---

create a new issue: upsilontools actor should if target message is of the right type (call/notification) right at the beginning (i would the the reply as well but we don't have the appropriate callback information... maybe instead of a reschan we should provide a link to the actor expecting the reply (as an alternate method))

```go

func (a *Actor) SendActor(msg *message.Message, res chan *message.Message) {
	if msg.CallbackMethod == nil {
		panic(fmt.Sprintf("[%s] Protocol violation: SendActor called with nil CallbackMethod for message %s. Use NotifyActor for fire-and-forget.", a.Name(), msg.TargetString()))
	}
	if res == nil {
		panic(fmt.Sprintf("[%s] Protocol violation: SendActor called with nil response channel for message %s. Provide a channel or use NotifyActor if you don't care about the response.", a.Name(), msg.TargetString()))
	}
	msg.ReplyChan = res
	a.Logger.WithFields(logrus.Fields{
		"message":      msg.String(),
		"message_type": msg.TargetString()}).Debug("Sending message")
	a.queue.Send(msg)
}

func (a *Actor) NotifyActor(msg *message.Message) {
	a.Logger.WithFields(logrus.Fields{
		"message":      msg.String(),
		"message_type": msg.TargetString()}).Debug("Notifying message")
	msg.ShouldBeRepliedTo = false
	a.queue.Send(msg)
}


```

---
```go

func (hc *HTTPController) forwardToWebhook(ctx actor.NotificationContext) {
	var action *api.ActionFeedback
	switch d := ctx.Msg.TargetMethod.(type) {
	case rulermethods.ControllerAttacked:
		action = &api.ActionFeedback{
			Type:     "attack",
			ActorID:  d.Attacker.ID.String(),
			TargetID: d.Entity.ID.String(),
			Damage:   d.Damage,
			PrevHP:   d.PrevHP,
			NewHP:    d.NewHP,
		}
	case rulermethods.ControllerMoved:
		action = &api.ActionFeedback{
			Type:    "move",
			ActorID: d.EntityID.String(),
			Path:    d.Path,
		}
	case rulermethods.ControllerPassed:
		action = &api.ActionFeedback{
			Type:    "pass",
			ActorID: d.EntityID.String(),
		}
	}

```

has no default to handle unexpected message (what about forfeit ? or simply no information)
at least add a default with a log.

---

```go

func (b *ArenaBridge) GetBoardState(matchID uuid.UUID, action *api.ActionFeedback) (api.BoardState, error) {
	b.mu.RLock()
	arena, ok := b.arenas[matchID]
	b.mu.RUnlock()
	if !ok {
		return api.BoardState{}, fmt.Errorf("arena %s not found", matchID)
	}

	res := make([]entity.Entity, 0, len(arena.Ruler.GameState.Entities))
    //vvvvv this is somehow crashing (on Next) I expect it's a race condition ??? but .... that doens't quite align with what happened so don't know. First fix everything else.
	for _, v := range arena.Ruler.GameState.Entities {
		res = append(res, v)
	}

	players, _ := arena.Metadata["Players"].([]api.Player)

	return api.NewBoardState(matchID, arena.Ruler.GameState.Grid, res, players, arena.Ruler.GameState.Turner.GetTurnState(), time.Now(), time.Now().Add(30*time.Second), arena.Ruler.GameState.WinnerTeamID, arena.Ruler.GameState.Version, action), nil
}
```

Clearly bypasses ruler's ownership of data at this step (there should have been an atd about this; once ruler.Start() has been called any access to ruler's gamestate is forbiden... Maybe we should gate GameState behind a Getter that panics if called once game started (need to ensure that ruler himself never calls it... heavily comment etc. and make gamestate private to the ruler)

---

So: why didn't this crash occurs with CLI first ? (or maybe it did and I wasn't aware?)

