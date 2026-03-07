package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusProcessController(c *gin.Context) {
	processId := c.Param("id")

	c.String(http.StatusOK, "Status Process: "+processId)
}
