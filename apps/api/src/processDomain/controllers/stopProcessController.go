package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StopProcessController(c *gin.Context) {
	processId := c.Param("id")

	c.String(http.StatusOK, "Stop Process: "+processId)
}
