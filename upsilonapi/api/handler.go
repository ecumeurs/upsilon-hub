package api

import (
	"net/http"

	"github.com/ecumeurs/upsilonapi/bridge"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/rulermethods"
	"github.com/ecumeurs/upsilontools/tools/messagequeue/message"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @spec-link [[api_go_battle_engine]]

// HandleArenaStart handles the start of a new arena; initializes a new ruler and returns the initial state.
func HandleArenaStart(c *gin.Context) {
	var req ArenaStartMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arenaID, state, err := bridge.Get().StartArena(req.Data.CallbackURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"arena_id":      arenaID,
		"initial_state": state,
	})
}

// HandleArenaAction handles an action in an arena; sends the action to the ruler.
func HandleArenaAction(c *gin.Context) {
	// extract StandardMessage first .

	idStr := c.Param("id")
	arenaID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid arena id"})
		return
	}

	var req ArenaActionMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r, ok := bridge.Get().GetArena(arenaID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "arena not found"})
		return
	}

	// Translate HTTP action to Ruler message
	// This is a simplified mapping; more logic needed for full support
	switch req.Data.Type {
	case "move":
		// Handle move...
		c.JSON(http.StatusOK, gin.H{"status": "accepted"})
	default:
		// Just notify the ruler for now with a generic message if type matches?
		// Better to implement specific methods
		r.NotifyActor(message.Create(nil, rulermethods.EndOfTurn{
			ControllerID: uuid.MustParse(req.Data.PlayerID),
			EntityID:     uuid.MustParse(req.Data.EntityID),
		}, nil))
		c.JSON(http.StatusOK, gin.H{"status": "accepted"})
	}
}
