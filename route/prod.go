package route

import (
	"scrapper/handler"

	"github.com/labstack/echo/v4"
)

func GeneralRoute(g *echo.Group) {
	g.GET("/scrap", func (c echo.Context) error {
		return handler.ScrapDenyaKolha(c)
	})
}	