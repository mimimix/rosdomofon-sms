package webServer

import (
	"context"
	"fmt"

	"domofon-api.gg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the request's origin matches the allowed origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Allow credentials

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// @title           ReCoin API
// @version         1.0

func New(config *config.Config, lc fx.Lifecycle) *gin.Engine {
	webServer := gin.Default()
	webServer.Use(gin.Recovery())

	webServer.Use(CORSMiddleware())

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					err := webServer.Run(fmt.Sprintf(":%d", config.HttpPort))
					if err != nil {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: nil,
		},
	)

	return webServer
}
