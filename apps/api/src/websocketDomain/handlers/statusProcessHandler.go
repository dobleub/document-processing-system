package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"nx-recipes/dps/lambda/interfaces"
	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	upgrader *websocket.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins, this must be changed in production
		},
	}
)

// @BasePath /ws

// StatusProcessHandler godoc
// @Summary Get All Process Status
// @Schemes
// @Description Get the current status of All Process connections via WebSocket
// @Tags Websocket Process
// @Accept json
// @Produce json
// @Success 200 {object} pd_interfaces.OperationListResponse "List of Processes"
// @Failure 400 {string} string "Bad Request"
// @Router /ws/status [get]
func StatusProcessHandler(c *gin.Context) {
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
	logger := baseLogger.With(zap.String("handler", "StatusProcessHandler"))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}
	defer func() {
		conn.Close()
		logger.Info("WebSocket connection closed")
	}()

	logger.Info("WebSocket connection established")

	// Set up ping/pong handlers for connection health
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start a goroutine to handle ping messages
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-pingTicker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					logger.Debug("Failed to send ping", zap.Error(err))
					cancel()
					return
				}
			}
		}
	}()

	// Track process snapshots to detect changes.
	// Snapshot format: <status>|<filesProcessedCount>
	allProcessSnapshot := make(map[string]string)
	updateTicker := time.NewTicker(time.Second)
	defer updateTicker.Stop()

	// Send initial state immediately
	if err := sendProcessUpdate(conn, state, allProcessSnapshot, logger); err != nil {
		logger.Error("Failed to send initial update", zap.Error(err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled, closing WebSocket")
			return

		case <-updateTicker.C:
			if err := sendProcessUpdate(conn, state, allProcessSnapshot, logger); err != nil {
				logger.Error("Failed to send update", zap.Error(err))
				return
			}
		}
	}
}

// sendProcessUpdate collects current process statuses and sends if there are changes
func sendProcessUpdate(conn *websocket.Conn, state *sync.Map, allProcessSnapshot map[string]string, logger *zap.Logger) error {
	processes := []pd_interfaces.OperationReview{}

	currentProcessIDs := make(map[string]bool)
	hasChanges := false

	state.Range(func(key, value interface{}) bool {
		if operationResponse, ok := value.(*pd_interfaces.OperationResponse); ok {
			opStatus := operationResponse.Status
			currentProcessIDs[opStatus.ID] = true

			filesProcessed := []string{}
			for _, file := range opStatus.Result.FilesProcessed {
				filesProcessed = append(filesProcessed, filepath.Base(file))
			}
			filesToProcess := []string{}
			for _, file := range opStatus.Result.FilesToProcess {
				filesToProcess = append(filesToProcess, filepath.Base(file))
			}

			processes = append(processes, pd_interfaces.OperationReview{
				ID:                  opStatus.ID,
				Status:              string(opStatus.Status),
				Error:               opStatus.Error,
				StartedAt:           opStatus.StartedAt,
				EstimatedCompletion: opStatus.EstimatedCompletion,
				FilesProcessed:      filesProcessed,
				FilesToProcess:      filesToProcess,
				CompletedAt:         opStatus.CompletedAt,
			})

			processSnapshot := string(opStatus.Status) + "|" + strconv.Itoa(len(filesProcessed)) + "|" + opStatus.CompletedAt

			// Track changes in status and FilesProcessed length.
			if existingSnapshot, exists := allProcessSnapshot[opStatus.ID]; !exists || existingSnapshot != processSnapshot {
				hasChanges = true
			}
			allProcessSnapshot[opStatus.ID] = processSnapshot
		}
		return true
	})

	// Clean up removed processes from tracking map to prevent memory leak
	for processID := range allProcessSnapshot {
		if !currentProcessIDs[processID] {
			delete(allProcessSnapshot, processID)
			hasChanges = true
		}
	}

	// Only send if there are changes
	if !hasChanges {
		return nil
	}

	jsonData, err := json.Marshal(processes)
	if err != nil {
		logger.Error("Failed to marshal process data", zap.Error(err))
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return err
	}

	logger.Debug("Sent process update", zap.Int("process_count", len(processes)))
	return nil
}
