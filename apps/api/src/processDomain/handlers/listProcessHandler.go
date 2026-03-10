package controllers

import (
	"net/http"
	"path/filepath"
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
// @Success 200 {object} pd_interfaces.OperationListResponse "List of Processes"
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

	processes := pd_interfaces.OperationListResponse{}
  processes.Initialize()
	state.Range(func(key, value interface{}) bool {
		if operationResponse, ok := value.(*pd_interfaces.OperationResponse); ok {
			opStatus := operationResponse.Status

			filesProcessed := []string{}
			for _, file := range opStatus.Result.FilesProcessed {
				filesProcessed = append(filesProcessed, filepath.Base(file))
			}
			filesToProcess := []string{}
			for _, file := range opStatus.Result.FilesToProcess {
				filesToProcess = append(filesToProcess, filepath.Base(file))
			}

			processes.AddProcess(pd_interfaces.OperationReview{
				ID:                  opStatus.ID,
				Status:              string(opStatus.Status),
				Error:               opStatus.Error,
				StartedAt:           opStatus.StartedAt,
				EstimatedCompletion: opStatus.EstimatedCompletion,
				FilesProcessed:      filesProcessed,
				FilesToProcess:      filesToProcess,
				CompletedAt:         opStatus.CompletedAt,
			})
		}
		return true
	})

  processes.OrderProcesses()
	duration := time.Since(start_time)
	logger.Info("List Processes", zap.Duration("duration", duration))

	c.JSON(http.StatusOK, gin.H{"processes": processes.Processes})
}
