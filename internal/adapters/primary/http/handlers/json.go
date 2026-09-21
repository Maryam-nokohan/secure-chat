package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func wantsJSON(c *gin.Context) bool {
	accept := c.GetHeader("Accept")
	return strings.Contains(accept, "application/json") ||
		c.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

func bindAuthBody(c *gin.Context, dest any) error {
	ct := c.ContentType()
	if strings.Contains(ct, "application/json") {
		return c.ShouldBindJSON(dest)
	}
	return c.ShouldBind(dest)
}
