package controllers

import (
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/genai"

	"nx-recipes/dps/lambda/interfaces"
	pd_interfaces "nx-recipes/dps/lambda/src/processDomain/interfaces"
	pd_lib "nx-recipes/dps/lambda/src/processDomain/lib"
)

// @BasePath /process

// StartProcessController godoc
// @Summary Start a Process
// @Schemes
// @Description Start a new process with the provided data
// @Tags Process
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "Process Data"
// @Success 200 {string} string "Process Started"
// @Failure 400 {string} string "Bad Request"
// @Router /process/start [post]
func StartProcessHandler(c *gin.Context) {
	start_time := time.Now()

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
	logger := baseLogger.With(zap.String("handler", "StartProcessHandler"), zap.Any("state", state))

	var mcpClient *genai.Client
	if mcpClientFromCtx, ok := c.Request.Context().Value(interfaces.McpClient).(*genai.Client); ok {
		mcpClient = mcpClientFromCtx
	} else {
		mcpClient = c.MustGet(string(interfaces.McpClient)).(*genai.Client)
	}

	// check if the request is POST
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	// process the files
	id := uuid.New().String()

	operationResponse := &pd_interfaces.OperationResponse{}
	operationResponse.Initialize(id)
	state.Store(id, operationResponse)

	// process files from a directory
	// [x] Use FileManager to list files from a specified directory
	// [x] Read the content of each file
	// [x] Generate summaries for each file
	// [x] Extract statistics from each file

	// get current path
	currentPath, err := os.Getwd()
	if err != nil {
		logger.Error("Error getting current path", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	FileProcessing := &pd_lib.FileProcessing{
		Path:      currentPath + "/targetFiles", // This should be configurable in a real application
		State:     state,
		McpClient: mcpClient,
		Log:       logger,
	}
	go FileProcessing.ProcessFilesFromDirectory(id) // Process files in a separate goroutine to avoid blocking the request

	duration := time.Since(start_time)
	logger.Info("Process Started", zap.Duration("duration", duration))

	c.JSON(http.StatusOK, gin.H{"message": "Process Started", "id": id})
}
