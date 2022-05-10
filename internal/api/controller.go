package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/logger"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/service"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}
type ServerParams struct {
	Cfg  *Config
	Log  *logrus.Entry
	Srvc *service.Service
}

type Controller struct {
	ServerInterface

	srvc *service.Service
}

func NewController(srvc *service.Service) *Controller {
	return &Controller{
		srvc: srvc,
	}
}

func NewServer(params *ServerParams) *http.Server {
	ctrl := NewController(params.Srvc)

	router := gin.New()
	router.Use(logger.NewLoggingMiddleware(params.Log)())

	if params.Cfg.Prefix == "" {
		params.Cfg.Prefix = "/"
	}

	RegisterHandlersWithOptions(router, ctrl, GinServerOptions{
		BaseURL: params.Cfg.Prefix,
	})

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
