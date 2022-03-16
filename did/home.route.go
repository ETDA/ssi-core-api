package did

import (
	"github.com/labstack/echo/v4"
	"ssi-gitlab.teda.th/ssi/core"
)

func NewHomeHTTPHandler(r *echo.Echo) {
	home := &HomeController{}

	r.GET("/", core.WithHTTPContext(home.Get))

}
