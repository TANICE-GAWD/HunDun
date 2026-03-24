package main

import (
	"github.com/gin-gonic/gin"
	"backend/internal/hub"
)

func main() {
	r := gin.Default()

	h := hub.NewHub()

	go hub.Run(h)
	go h.StartTick()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "StatusOK")
	})

	r.GET("/ws", func(c *gin.Context) {
		hub.ServeWs(h, c)
	})

	r.Run(":8080")
}