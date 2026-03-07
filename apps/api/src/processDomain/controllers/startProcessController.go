package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartProcessController(c *gin.Context) {
	c.String(http.StatusOK, "Start Process")
}
