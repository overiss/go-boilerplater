package scaffold

import "fmt"

const (
	creatorTemplate = `package app

import (
	"context"
	"fmt"
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

func (a *Application) failInit(ctx context.Context, component string, err error) {
	_ = ctx
	log.Printf("cannot init %s: %v", component, err)
	panic(fmt.Errorf("cannot init %s: %w", component, err))
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
	configTemplate = `package config

type Application struct {
	Version string
}
`
	modelTemplate      = "package model\n"
	httpServerTemplate = `package http

type HTTP struct{}
`
	handlerContainerTemplate = `package handler

type Container struct{}
`
	serviceContainerTemplate = `package service

type Container struct{}
`
	varsTemplate = `package vars
`
	utilsTemplate = `package utils
`
)

func mainGoTemplate(serviceName string) string {
	return fmt.Sprintf(`package main

import (
	"context"
	"os/signal"
	"syscall"

	"log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("starting service: %s")
	<-ctx.Done()
	log.Printf("stopping service: %s")
}
`, serviceName, serviceName)
}

func appTemplate(moduleName string) string {
	return fmt.Sprintf(`package app

import (
	"context"
	"log"

	"%s/internal/behavior"
	"%s/internal/config"
	"%s/internal/server"
	httpHandler "%s/internal/server/http/handler"
	"%s/internal/service"
)

type Application struct {
	cfg              *config.Application
	netContainer     *server.Container
	handlerContainer *httpHandler.Container
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

func serverContainerTemplate(moduleName string) string {
	return fmt.Sprintf(`package server

import serverHttp "%s/internal/server/http"

type Container struct {
	privateHTTP *serverHttp.HTTP
	publicHTTP  *serverHttp.HTTP
}

func NewContainer(privateHTTP, publicHTTP *serverHttp.HTTP) *Container {
	return &Container{
		privateHTTP: privateHTTP,
		publicHTTP:  publicHTTP,
	}
}

func (c *Container) PrivateHTTP() *serverHttp.HTTP {
	return c.privateHTTP
}

func (c *Container) PublicHTTP() *serverHttp.HTTP {
	return c.publicHTTP
}
`, moduleName)
}
