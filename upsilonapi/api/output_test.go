package api

import (
	"testing"
	"time"

	"github.com/ecumeurs/upsilonbattle/battlearena/entity"
	"github.com/ecumeurs/upsilonbattle/battlearena/ruler/turner"
	"github.com/ecumeurs/upsilonmapdata/grid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBoardStateWinnerID(t *testing.T) {
	matchID := uuid.New()
	g := grid.NewGrid(10, 10, 1)
	entities := []entity.Entity{}
	players := []Player{}
	ts := turner.TurnState{}
	startTime := time.Now()
	timeout := time.Now().Add(30 * time.Second)
	winnerID := uuid.New()

	// Test with a winner
	bs := NewBoardState(matchID, g, entities, players, ts, startTime, timeout, winnerID)
	assert.Equal(t, winnerID.String(), bs.WinnerID, "WinnerID should be populated in BoardState")

	// Test without a winner (uuid.Nil)
	bs = NewBoardState(matchID, g, entities, players, ts, startTime, timeout, uuid.Nil)
	assert.Equal(t, "", bs.WinnerID, "WinnerID should be empty when uuid.Nil is passed")
}
