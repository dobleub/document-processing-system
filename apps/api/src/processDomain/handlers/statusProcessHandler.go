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

// StatusProcessController godoc
// @Summary Get Process Status
// @Schemes
// @Description Get the status of a process by its ID
// @Tags Process
// @Accept json
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {object} pd_interfaces.OperationStatus "Process Status"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Process Not Found"
// @Router /process/status/{id} [get]
func StatusProcessHandler(c *gin.Context) {
	state := c.MustGet(string(interfaces.StateKey)).(*sync.Map)
	logger := c.MustGet(string(interfaces.LoggerKey)).(*zap.Logger).With(zap.String("handler", "StartProcessHandler"), zap.Any("state", state))
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
	logger.Info("Status Process", zap.String("processId", processId), zap.Duration("duration", duration))

	if val, ok := state.Load(processId); ok {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, val.(*pd_interfaces.OperationStatus))
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
}
