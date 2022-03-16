package did

import (
	"github.com/labstack/echo/v4"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/middlewares"
)

func NewVCHTTPHandler(r *echo.Echo) {
	vc := &VCController{}
	r.GET("/vc/:id", core.WithHTTPContext(vc.Find))
	r.POST("/vc", core.WithHTTPContext(vc.Register), middlewares.ValidateSignatureMessageMiddleware)
}
