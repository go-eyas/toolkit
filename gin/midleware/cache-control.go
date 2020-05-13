package middleware

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

var StaticRegexp = regexp.MustCompile("/*/*.(js|css|png|jpg|woff|tff|oet)")

func CacheControl(filter *regexp.Regexp) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.RequestURI
		if filter != nil && filter.MatchString(path) {
			c.Header("Cache-Control", "public, max-age=31536000")
		}

		c.Next()
		status := c.Writer.Status()
		if status > 300 || status < 200 {
			c.Header("Cache-Control", "")
		}
	}
}
