package controllers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"nx-recipes/dps/lambda/interfaces"
	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"
)

// @BasePath /process

// ListProcessController godoc
// @Summary List Processes
// @Schemes
// @Description List all processes with their current status
// @Tags Process
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "Process Data"
// @Success 200 {string} string "Process List"
// @Failure 400 {string} string "Bad Request"
// @Router /process/list [get]
func ListProcessHandler(c *gin.Context) {
	var state *sync.Map
	if stateFromCtx, ok := c.Request.Context().Value(interfaces.StateKey).(*sync.Map); ok {
		state = stateFromCtx
	} else {
		state = c.MustGet(string(interfaces.StateKey)).(*sync.Map)
	}

	var baseLogger *zap.Logger
	if loggerFromCtx, ok := c.Request.Context().Value(interfaces.LoggerKey).(*zap.Logger); ok {
		baseLogger = loggerFromCtx
	} else {
		baseLogger = c.MustGet(string(interfaces.LoggerKey)).(*zap.Logger)
	}
	logger := baseLogger.With(zap.String("handler", "ListProcessHandler"), zap.Any("state", state))
	start_time := time.Now()
	// check if the request is GET
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	var processes map[string]map[string]interface{} = make(map[string]map[string]interface{})
	state.Range(func(key, value interface{}) bool {
		if operationResponse, ok := value.(*pd_interfaces.OperationResponse); ok {
			opStatus := operationResponse.Status
			processes[opStatus.ID] = map[string]interface{}{
				"process_id":           opStatus.ID,
				"status":               opStatus.Status,
				"error":                opStatus.Error,
				"started_at":           opStatus.StartedAt,
				"estimated_completion": opStatus.EstimatedCompletion,
			}
		}
		return true
	})

	duration := time.Since(start_time)
	logger.Info("List Processes", zap.Duration("duration", duration))

	c.JSON(http.StatusOK, gin.H{"processes": processes})
}
