package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListProcessHandler(c *gin.Context) {
	c.String(http.StatusOK, "List all Process")
}
