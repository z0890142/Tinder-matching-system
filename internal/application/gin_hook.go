package app

import (
	"tinderMatchingSystem/internal/middleware"
	"tinderMatchingSystem/internal/web/restful"

	"github.com/gin-gonic/gin"
)

func initGinApplicationHook(app *Application) error {
	if app.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.EnableJsonDecoderUseNumber()
	g := gin.New()

	g.Use(gin.Recovery())
	g.Use(middleware.ErrorHandler())

	router := g.Group("/tinder_system/v1")
	restfulHandler := restful.NewRestfulHandler()

	router.POST("/persons", restfulHandler.AddSinglePersonAndMatch)
	router.DELETE("/persons/:name", restfulHandler.RemoveSinglePerson)
	router.GET("/persons", restfulHandler.QuerySinglePeople)

	app.srv.Handler = g
	app.Logger.Debug("init gin http server")

	return nil
}
