package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const loggerKey = "logger"

type loggerCtxKey struct{}

func DefaultLogger() *logrus.Entry {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{})

	return logrus.NewEntry(log)
}

func GetLogger(c *gin.Context) *logrus.Entry {
	le, ok := c.Get(loggerKey)
	if !ok {
		return DefaultLogger()
	}

	entry, ok := le.(*logrus.Entry)
	if !ok {
		return DefaultLogger()
	}

	return entry
}

func NewLoggingMiddleware(logEntry *logrus.Entry) func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			nextEntry := logEntry.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			})

			c.Set(loggerKey, nextEntry)
		}
	}
}
