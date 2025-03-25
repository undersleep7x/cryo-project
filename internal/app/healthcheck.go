package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// handle /ping route call and return ok to confirm healthy service
var Ping = func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PONG"})
}
