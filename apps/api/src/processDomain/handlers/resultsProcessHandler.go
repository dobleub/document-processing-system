package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResultsProcessHandler(c *gin.Context) {
	c.String(http.StatusOK, "Result Process")
}
