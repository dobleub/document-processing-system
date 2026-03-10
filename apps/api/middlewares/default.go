package middlewares

import (
	"context"
	"net/http"
	"sync"

	"nx-recipes/dps/lambda/config"
	"nx-recipes/dps/lambda/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"google.golang.org/genai"
)

func Setup(router *gin.Engine, appLogger *zap.Logger, env *config.Config, dbContext interfaces.MongoDBContext, state *sync.Map, mcpClient *genai.Client) {
	// Add logging and recovery middleware
	gin.DefaultWriter = colorable.NewColorableStdout()
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))
	router.Use(gin.Recovery())
	// Add custom context for each request
	router.Use(func(c *gin.Context) {
		reqCtx := c.Request.Context()
		reqCtx = context.WithValue(reqCtx, interfaces.LoggerKey, appLogger)
		reqCtx = context.WithValue(reqCtx, interfaces.ConfigKey, env)
		reqCtx = context.WithValue(reqCtx, interfaces.MongodbKey, dbContext)
		reqCtx = context.WithValue(reqCtx, interfaces.StateKey, state)
		reqCtx = context.WithValue(reqCtx, interfaces.McpClient, mcpClient)
		c.Request = c.Request.WithContext(reqCtx)
		c.Next()
	})
	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	// Add custom error handling middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				appLogger.Error("Request error", zap.Error(e.Err))
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
	})
	// Add basic auth middleware for protected routes
	router.Use(func(c *gin.Context) {
		// Allow unauthenticated access to health. Swagger and docs endpoints are protected but for simplicity we will allow unauthenticated access to them as well
		if c.FullPath() == "/health" || c.FullPath() == "/swagger/*any" {
			c.Next()
			return
		}

		// Check for Authorization header first
		var token string
		authHeader := c.GetHeader("Authorization")
    connectionType := c.GetHeader("Upgrade")

    if connectionType == "websocket" {
      appLogger.Info("WebSocket connection detected, using query parameter for authentication")
      // Fallback to query parameter for WebSocket connections
      token = c.Query("token")
    } else {
      appLogger.Info("Standard HTTP connection detected, using Authorization header for authentication")
      // Header-based authentication (REST API)
      token = authHeader[len("Bearer "):]
    }

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if token != env.APIAuthToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	})
}
