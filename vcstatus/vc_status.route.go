package vcstatus

import (
	"github.com/labstack/echo/v4"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/middlewares"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/views"
)

func NewVCStatusHTTPHandler(r *echo.Echo) {
	vcStatus := &VCStatusController{}
	r.POST("/vp/verify", core.WithHTTPContext(vcStatus.VerifyVP), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}), middlewares.ValidateSignatureJWTMessageMiddleware)
	r.POST("/vc/verify", core.WithHTTPContext(vcStatus.VerifyVC), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}), middlewares.ValidateSignatureJWTMessageMiddleware)
	r.POST("/vc/status", core.WithHTTPContext(vcStatus.Add), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}), middlewares.ValidateSignatureMessageMiddleware)
	r.GET("/vc/status", core.WithHTTPContext(vcStatus.All), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}))
	r.GET("/vc/status/:id", core.WithHTTPContext(vcStatus.Find), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}))
	r.PUT("/vc/status/:id", core.WithHTTPContext(vcStatus.Update), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}), middlewares.ValidateSignatureMessageMiddleware)
	r.POST("/vc/status/tags", core.WithHTTPContext(vcStatus.Tag), core.MockMiddleware(&views.VCStatus{}, &core.MockMiddlewareOptions{
		IsDisabled: true,
	}), middlewares.ValidateSignatureMessageMiddleware)
}
