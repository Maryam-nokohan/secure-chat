package spa

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/maryam-nokohan/secure-chat/pkg"
)

func Register(r *gin.Engine, distDir string) {
	index := filepath.Join(distDir, "index.html")
	if _, err := os.Stat(index); err != nil {
		pkg.LogInfo("React build not found at " + index + " - run `npm run build` in ./frontend; /login and /register will 404 until then")
	}

	r.Static("/assets", filepath.Join(distDir, "assets"))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		isRead := c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead
		if !isRead || strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/assets/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		file := filepath.Join(distDir, filepath.Clean("/"+path))
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			c.File(file)
			return
		}

		c.Header("Cache-Control", "no-cache")
		c.File(index)
	})
}
