package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Router /process/upload [post]
func UploadProcessHandler(c *gin.Context) {
	start_time := time.Now()
	// check if the request is POST
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}
	// check if the request contains a file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	// if not check if the request contains a list of files
	files, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Files are required"})
		return
	}
	if len(files.File["files"]) == 0 && file.Filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Files are required"})
		return
	}
	// process the files
	id := uuid.New()
	c.JSON(http.StatusOK, gin.H{"message": "Process Started", "id": id.String()})

	duration := time.Since(start_time)
	// log the duration of the process
	c.JSON(http.StatusOK, gin.H{"message": "Process Completed", "duration": duration.String()})
}
