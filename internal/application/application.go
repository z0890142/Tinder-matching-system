package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"tinderMatchingSystem/internal/c"

	"tinderMatchingSystem/config"

	"go.uber.org/zap"
)

var defaultApplication *Application

type Application struct {
	mu     sync.Mutex
	Logger *zap.SugaredLogger
	Config *config.Config
	srv    *http.Server
	// Init and destroy hooks
	initHooks    []ApplicationHook
	destroyHooks []ApplicationHook
}

type ApplicationHook func(*Application) error

// Default get default application
func Default() *Application {
	if defaultApplication != nil {
		return defaultApplication
	}
	o := &Application{
		Config:       config.GetConfig(),
		initHooks:    make([]ApplicationHook, 0),
		destroyHooks: make([]ApplicationHook, 0),
		srv:          &http.Server{},
	}

	initSlice := []ApplicationHook{initLoggerApplicationHook}
	for _, hook := range initSlice {
		err := hook(o)
		if err != nil {
			panic(err)
		}
	}
	o.AddInitHook(initGinApplicationHook)

	defaultApplication = o

	return o
}

// Run application
func (app *Application) Run() {
	app.callInitHooks()

	errc := make(chan error)

	go func() {
		addr := fmt.Sprintf("%s:%d", app.Config.Service.Host, app.Config.Service.Port)
		if app.Logger != nil {
			app.Logger.Info("running server on: ", addr)
		}
		app.srv.Addr = fmt.Sprintf("%s:%d", app.Config.Service.Host, app.Config.Service.Port)
		errc <- app.srv.ListenAndServe()
	}()

	app.Logger.Error((fmt.Sprintf("application run error: %s", <-errc)))
}

// AddInitHook add init callback function
func (app *Application) AddInitHook(f ApplicationHook) {
	app.initHooks = append(app.initHooks, f)
}

func (app *Application) AddDestroyHook(f ApplicationHook) {
	app.destroyHooks = append(app.destroyHooks, f)
}

// Shutdown shundown service
func (app *Application) Shutdown() {
	app.mu.Lock()
	defer app.mu.Unlock()

	if app.Logger != nil {
		app.Logger.Warn("shutdowning")
	}

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.srv.Shutdown(c); err != nil {
		app.Logger.Error("srv.Shutdown:", err)
	}
	select {
	case <-c.Done():
		app.Logger.Info("Graceful Shutdown http server")
		done := make(chan struct{})
		go app.callDestroyHooks(done)
		t := time.NewTimer(5 * time.Second)

		select {
		case <-done:
			break
		case <-t.C:
			if app.Logger != nil {
				app.Logger.Warn("timeout: application destroy hooks interrupted")
			}
			break
		}
	}

}

func (app *Application) callInitHooks() {
	for _, hook := range app.initHooks {
		if err := hook(app); err != nil {
			panic(err)
		}
	}
}

func (app *Application) callDestroyHooks(done chan struct{}) {
	for i := len(app.destroyHooks); i > 0; i-- {
		hook := app.destroyHooks[i-1]
		if err := hook(app); err != nil {
			app.Logger.Error("calling application destroy hook error: ", err.Error())
		}
	}

	done <- struct{}{}
}

func (app *Application) IsProduction() bool {
	return app.Environment() == c.EnvProduction
}

// Environment get Environment
func (app *Application) Environment() string {
	return strings.ToLower(app.Config.Env)
}
