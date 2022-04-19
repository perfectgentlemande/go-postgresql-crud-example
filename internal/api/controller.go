package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/logger"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}
type ServerParams struct {
	Cfg *Config
	Log *logrus.Entry
}

type Controller struct {
	ServerInterface
}

func NewController() *Controller {
	return &Controller{}
}

func NewServer(params *ServerParams) *http.Server {
	ctrl := NewController()

	router := gin.New()
	router.Use(logger.NewLoggingMiddleware(params.Log)())
	apiRouter := gin.New()
	RegisterHandlers(apiRouter, ctrl)

	if params.Cfg.Prefix == "" {
		params.Cfg.Prefix = "/"
	}
	router.Group(params.Cfg.Prefix, apiRouter.Handlers...)

	return &http.Server{
		Addr:    params.Cfg.Addr,
		Handler: router,
	}
}

func WriteError(c *gin.Context, status int, message string) {
	c.JSON(status, APIError{Message: message})
}
func WriteSuccessful(c *gin.Context, payload interface{}) {
	c.JSON(http.StatusOK, payload)
}
