package main

import (
	"github.com/ecumeurs/upsilonapi/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	logrus.Info("Starting UpsilonAPI server on :8081")

	// Internal Arena Management
	internal := r.Group("/internal")
	{
		internal.POST("/arena/start", handler.HandleArenaStart)
		internal.POST("/arena/:id/action", handler.HandleArenaAction)
	}

	// V1 API
	v1 := r.Group("/v1")
	{
		v1.GET("/match/stats/active", handler.HandleGetActiveMatchStats)
	}

	if err := r.Run(":8081"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
