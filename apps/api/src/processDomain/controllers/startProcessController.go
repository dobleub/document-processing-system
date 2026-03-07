package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
func StartProcessController(c *gin.Context) {
	c.String(http.StatusOK, "Start Process")
}
