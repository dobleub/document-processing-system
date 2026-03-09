package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/genai"

	"nx-recipes/dps/lambda/interfaces"
	summarizerDomainLib "nx-recipes/dps/lambda/src/summarizerDomain/lib"
)

// @BasePath /summarizer

// httpHandler godoc
// @Summary http Handler for Summarizer
// @Schemes
// @Description Handle http requests for summarization
// @Tags Summarizer
// @Accept json
// @Produce json
// @Param data body sd_interfaces.SummarizeInput true "Summarize Input"
// @Success 200 {string} string "Summarize Output"
// @Failure 400 {string} string "Bad Request"
// @Router /summarizer/summarize [post]
func HttpSummarizerHandler(c *gin.Context) {
	var baseLogger *zap.Logger
	if loggerFromCtx, ok := c.Request.Context().Value(interfaces.LoggerKey).(*zap.Logger); ok {
		baseLogger = loggerFromCtx
	} else {
		baseLogger = c.MustGet(string(interfaces.LoggerKey)).(*zap.Logger)
	}
	logger := baseLogger.With(zap.String("handler", "HttpSummarizerHandler"))

	mcpClient := c.Request.Context().Value(interfaces.McpClient).(*genai.Client)

	start_time := time.Now()
	// check if the request is POST
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	// check if the request contains a single file
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		logger.Error("Failed to get file from request", zap.Error(err))
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
	}

	// check size of the file(s) not exceed ~32MB
	const maxFileSize = 32 << 20 // 32MB

	// process the files
	content := []byte{}
	fileName := ""
	fileSize := int64(0)
	summary := map[string]interface{}{}
	// process single file if exist
	if file != nil {
		fileName = fileHeader.Filename
		fileSize = fileHeader.Size
		if fileSize > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the limit of 32MB"})
			return
		}
		logger.Info("Processing single file", zap.Int("fileCount", 1), zap.String("singleFile", fileName))
		content, err = io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
	}
	if len(body) > 0 {
		if len(body) > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body size exceeds the limit of 32MB"})
			return
		}
		logger.Info("Processing request body", zap.Int("bodySize", len(body)))
		content = body
	}

	summaryContent, err := summarizerDomainLib.SummarizeContent(c, mcpClient, string(content), 100)
	if err != nil {
		logger.Error("Failed to summarize content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to summarize content"})
		return
	}
	// append summary to the list
	summary = map[string]interface{}{
		"fileName": fileName,
		"fileSize": fileSize,
		"summary":  summaryContent,
	}
	duration := time.Since(start_time)
	// log the duration of the process
	c.JSON(http.StatusOK, gin.H{"message": "Process Completed", "duration": duration.String(), "summary": summary})
}
