package main

import (
	"scrapper/route"
	"scrapper/sigs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/color"
)
func injectCors(e *echo.Echo) {
	devMode := true
	if (devMode) {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		  }))

	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		  }))
	}
}

func main() {


	println(color.Red(sigs.Signature))

	e := echo.New()

	injectCors(e)

	group := e.Group("/api")

	route.GeneralRoute(group)

	e.Logger.Fatal(e.Start(":8000"))
	
}
