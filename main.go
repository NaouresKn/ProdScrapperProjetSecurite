package main

import (
	"scrapper/route"
	"scrapper/sigs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
)

func main() {


	println(color.Red(sigs.Signature))

	e := echo.New()

	group := e.Group("/api")

	route.GeneralRoute(group)

	e.Logger.Fatal(e.Start(":8000"))
	
}
