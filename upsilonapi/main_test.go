package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	internal := r.Group("/internal")
	{
		internal.POST("/arena/start", handleArenaStart)
		internal.POST("/arena/:id/action", handleArenaAction)
	}
	return r
}

func TestArenaStartEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(ArenaStartRequest{
		MatchID:     "test-match",
		CallbackURL: "http://localhost:9999/webhook",
	})
	req, _ := http.NewRequest("POST", "/internal/arena/start", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "arena_id")
	assert.Contains(t, resp, "initial_state")
}
