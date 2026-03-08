package middlewares

import (
	"net/http"
	"sync"

	"nx-recipes/dps/lambda/config"
	"nx-recipes/dps/lambda/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
)

func Setup(router *gin.Engine, appLogger *zap.Logger, env *config.Config, dbContext interfaces.MongoDBContext, state *sync.Map) {
	// Add logging and recovery middleware
	gin.DefaultWriter = colorable.NewColorableStdout()
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))
	router.Use(gin.Recovery())
	// Add custom context for each request
	router.Use(func(c *gin.Context) {
		// Add logger to context for each request
		c.Set(string(interfaces.LoggerKey), appLogger)
		// Add appContext to context for each request
		c.Set(string(interfaces.ConfigKey), env)
		c.Set(string(interfaces.MongodbKey), dbContext)
    // Add state to context for each request
    c.Set(string(interfaces.StateKey), state)
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
		if c.FullPath() == "/health" || c.FullPath() == "/swagger/*any" {
			c.Next()
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		token := authHeader[len("Bearer "):]
		if token != env.APIAuthToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	})
}
