package main

import (
	"github.com/ecumeurs/upsilonapi/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	logrus.Info("Starting UpsilonAPI server on :8080")

	// Internal Arena Management
	internal := r.Group("/internal")
	{
		internal.POST("/arena/start", api.HandleArenaStart)
		internal.POST("/arena/:id/action", api.HandleArenaAction)
	}

	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
