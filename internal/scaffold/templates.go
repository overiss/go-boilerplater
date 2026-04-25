package scaffold

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	creatorTemplate = `package app

import (
	"context"
	"log"
)

func (a *Application) initNetContainer() {
	// TODO: initialize network container and append starters.
}

func (a *Application) initHandlerContainer(ctx context.Context) {
	if a.serviceContainer == nil {
		log.Printf("serviceContainer is not initialized. use initServiceContainer before initHandlerContainer")
		return
	}

	_ = ctx
	// TODO: initialize transport handlers from a.serviceContainer.
}

func (a *Application) initServiceContainer(ctx context.Context) {
	_ = ctx
	// TODO: initialize service container and append starters/readiness checks.
}
`
	routingTemplate = `package app

import (
	"context"
	"net/http"
	"strings"
)

func (a *Application) expose(ctx context.Context) {
	if a.handlerContainer == nil {
		return
	}
	a.publicRoutes()
	a.privateRoutes(ctx)
}

func (a *Application) publicRoutes() {
	// TODO: register public routes.
}

func (a *Application) privateRoutes(ctx context.Context) {
	_ = ctx
	_ = strings.TrimSpace(http.MethodGet)
	// TODO: register private routes and readiness checks.
}
`
	behaviorTemplate = `package behavior

import "context"

type Starter interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type Readiness interface {
	Name() string
	IsReady() bool
}
`
	modelTemplate            = "package model\n"
	handlerContainerTemplate = `package hanlderHttp

type Container struct{}
`
	serviceContainerTemplate = `package service

type Container struct{}
`
	varsTemplate = `package vars
`
)

func mainGoTemplate(moduleName string) string {
	return fmt.Sprintf(`package main

import (
	"context"
	"log"

	"%s/internal/app"
	"%s/internal/config"
	"%s/pkg/utils"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.Init()
	erg, ctx := errgroup.WithContext(context.Background())

	erg.Go(func() error {
		return utils.Listen(ctx)
	})

	log.Printf("config initialized")

	application := app.Init(ctx, cfg)
	log.Printf("service initialized successfully")
	log.Printf("IMPORTANT: replace all log.Printf/log.Println/log.Fatalf calls with your project logger")
	go application.Start(ctx)

	if err := erg.Wait(); err != nil {
		log.Printf("stopping application cause: %%v", err)
	}

	application.Stop(ctx)
	log.Printf("application stopped gracefully")
}
`, moduleName, moduleName, moduleName)
}

func appTemplate(moduleName string) string {
	return fmt.Sprintf(`package app

import (
	"context"
	"log"

	"%s/internal/behavior"
	"%s/internal/config"
	"%s/internal/server"
	hanlderHttp "%s/internal/server/http/handler"
	"%s/internal/service"
)

type Application struct {
	cfg              *config.Application
	netContainer     *server.Container
	handlerContainer *hanlderHttp.Container
	serviceContainer *service.Container
	starter          []behavior.Starter
	readiness        []behavior.Readiness
}

func Init(ctx context.Context, cfg *config.Application) *Application {
	app := &Application{
		cfg:       cfg,
		starter:   make([]behavior.Starter, 0),
		readiness: make([]behavior.Readiness, 0),
	}

	app.initNetContainer()
	app.initServiceContainer(ctx)
	app.initHandlerContainer(ctx)
	app.expose(ctx)

	log.Printf("application initialized successfully")
	log.Printf("todo: replace standard log package with your project logger")
	return app
}

func (a *Application) Start(ctx context.Context) {
	for _, starter := range a.starter {
		go starter.Start(ctx)
	}
}

func (a *Application) Stop(ctx context.Context) {
	for _, starter := range a.starter {
		starter.Stop(ctx)
	}
}
`, moduleName, moduleName, moduleName, moduleName, moduleName)
}

func configTemplate(serviceName string) string {
	configVarName := serviceVariableName(serviceName)

	return fmt.Sprintf(`package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Version = "v1.0.0" // would be rewritten on build stage

type (
	Application struct {
		HttpPrivateServer *HttpServer `+"`"+`envconfig:"HTTP_PRIVATE_SERVER" required:"true"`+"`"+`
		HttpPublicServer  *HttpServer `+"`"+`envconfig:"HTTP_PUBLIC_SERVER" required:"true"`+"`"+`

		Version string
		IsDebug bool `+"`"+`envconfig:"IS_DEBUG"`+"`"+`
	}

	HttpServer struct {
		Name         string        `+"`"+`envconfig:"NAME" required:"true"`+"`"+`
		Port         string        `+"`"+`envconfig:"PORT" required:"true"`+"`"+`
		ReadTimeout  time.Duration `+"`"+`envconfig:"READ_TIMEOUT" default:"15s"`+"`"+`
		WriteTimeout time.Duration `+"`"+`envconfig:"WRITE_TIMEOUT" default:"30s"`+"`"+`
		IdleTimeout  time.Duration `+"`"+`envconfig:"IDLE_TIMEOUT" default:"10s"`+"`"+`
	}
)

func Init() *Application {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: error loading .env file")
	}

	%s := new(Application)
	if err := envconfig.Process("", %s); err != nil {
		log.Fatalf("cannot process config for %s: %%v", err)
	}

	%s.Version = Version

	return %s
}
`, configVarName, configVarName, serviceName, configVarName, configVarName)
}

func serverContainerTemplate(moduleName string) string {
	return fmt.Sprintf(`package server

import serverHttp "%s/internal/server/http"

type Container struct {
	privateHttp *serverHttp.Http
	publicHttp  *serverHttp.Http
}

func NewContainer(privateHttp, publicHttp *serverHttp.Http) *Container {
	return &Container{
		privateHttp: privateHttp,
		publicHttp:  publicHttp,
	}
}

func (c *Container) PrivateHttp() *serverHttp.Http {
	return c.privateHttp
}

func (c *Container) PublicHttp() *serverHttp.Http {
	return c.publicHttp
}
`, moduleName)
}

func httpServerTemplate(moduleName string) string {
	return fmt.Sprintf(`package serverHttp

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"%s/internal/config"
)

type Http struct {
	server *http.Server
	router *gin.Engine
	name   string
}

func NewHttp(cfg *config.HttpServer, isDebug bool) *Http {
	if isDebug {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Http{
		name:   cfg.Name,
		server: server,
		router: router,
	}
}

func (s *Http) Start(_ context.Context) {
	log.Printf("%%s starting", s.name)

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("%%s failed to start: %%v", s.name, err)
	}
}

func (s *Http) Stop(ctx context.Context) {
	log.Printf("shutting down server %%s", s.name)

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("error shutting down %%s: %%v", s.name, err)
	}
	log.Printf("%%s stopped successfully", s.name)
}

func (s *Http) Router() *gin.Engine {
	return s.router
}
`, moduleName)
}

const utilsTemplate = `package utils

import (
	"context"
	"os/signal"
	"syscall"
)

func Listen(ctx context.Context) error {
	sigCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()
	if sigCtx.Err() == context.Canceled && ctx.Err() == nil {
		return nil
	}
	if sigCtx.Err() == context.Canceled {
		return nil
	}
	return sigCtx.Err()
}
`

func serviceVariableName(serviceName string) string {
	parts := strings.FieldsFunc(serviceName, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	if len(parts) == 0 {
		return "app"
	}

	for i := range parts {
		parts[i] = strings.ToLower(parts[i])
	}

	varName := parts[0]
	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		varName += strings.ToUpper(p[:1]) + p[1:]
	}
	if varName == "" {
		return "app"
	}
	if unicode.IsDigit(rune(varName[0])) {
		return "svc" + varName
	}
	return varName
}
