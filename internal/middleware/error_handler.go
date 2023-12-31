package middleware

import (
	"net/http"
	"tinderMatchingSystem/internal/error"

	"tinderMatchingSystem/pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			logger.LoadExtra(map[string]interface{}{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"error":  c.Errors.String(),
			}).Error("API Error")
			switch e := err.Err.(type) {
			case error.Http:
				c.AbortWithStatusJSON(e.Detail.StatusCode, e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": "Service Unavailable"})
			}
		}
	}
}
