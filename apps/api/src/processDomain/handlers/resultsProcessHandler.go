package controllers

import (
	"net/http"
	"sync"
	"time"

	"nx-recipes/dps/lambda/interfaces"
	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @BasePath /process

// ResultsProcessController godoc
// @Summary Get Process Results
// @Schemes
// @Description Get the results of a process by its ID
// @Tags Process
// @Accept json
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {object} pd_interfaces.OperationStatus "Process Results"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Process Not Found"
// @Router /process/results/{id} [get]
func ResultsProcessHandler(c *gin.Context) {
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
	logger := baseLogger.With(zap.String("handler", "StatusProcessHandler"), zap.Any("state", state))

	start_time := time.Now()
	// check if the request is GET
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	processId := c.Param("id")
	// check if the processId is empty
	if processId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Process ID is required"})
		return
	}

	duration := time.Since(start_time)
	logger.Info("Analyze Process", zap.String("processId", processId), zap.Duration("duration", duration))

	if val, ok := state.Load(processId); ok {
		opResponse := val.(*pd_interfaces.OperationResponse)
		opAnalysis := map[string]interface{}{
			"progress": opResponse.Status.Progress,
			"analysis": opResponse.Analysis,
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, opAnalysis)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
}
