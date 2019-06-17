package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConf struct {
	GinMode        string        `mapstructure:"gin_mode"`
	Address        string        `mapstructure:"address"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes uint          `mapstructure:"max_header_bytes"`
}

type Controller interface {
	Actions() []Action
}

type Action struct {
	HttpMethod   string
	RelativePath string
	ActionExec   func(ctx *gin.Context)
}

type Server interface {
	Engine() *gin.Engine
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	srv    *http.Server
	engine *gin.Engine
}

func (s *server) Engine() *gin.Engine {
	return s.engine
}

func (s *server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func NewServer(
	controller Controller,
	config *ServerConf,
) Server {

	s := &server{}
	s.setup(controller, config)

	return s
}

func (s *server) setup(controller Controller, config *ServerConf) {

	gin.SetMode(config.GinMode)

	s.engine = gin.New()
	s.engine.Use(gin.Recovery())

	if config.GinMode != gin.ReleaseMode {
		s.engine.Use(gin.Logger())
	}

	for _, action := range controller.Actions() {
		a := action
		s.engine.Handle(a.HttpMethod, a.RelativePath, a.ActionExec)
	}

	s.srv = &http.Server{
		Addr:           config.Address,
		Handler:        s.engine,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: int(config.MaxHeaderBytes),
	}
}
