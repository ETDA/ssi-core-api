package did

import (
	"github.com/labstack/echo/v4"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/middlewares"
	"ssi-gitlab.teda.th/ssi/core"
)

func NewDIDHTTPHandler(r *echo.Echo) {
	did := &DIDController{}
	key := &KeyController{}

	r.POST("/did", core.WithHTTPContext(did.Register), middlewares.ValidateSignatureMessageMiddleware)
	r.POST("/did/:did/recoverer/register", core.WithHTTPContext(did.AddRecoverer), middlewares.ValidateSignatureMessageMiddleware)
	r.GET("/did/:did/document/latest", core.WithHTTPContext(did.Find))
	r.GET("/did/:did/document/history", core.WithHTTPContext(did.History))
	r.GET("/did/:did/document/:version_id", core.WithHTTPContext(did.FindByVersion))
	r.GET("/did/:did/nonce", core.WithHTTPContext(did.GetNonce))
	r.POST("/did/:did/keys", core.WithHTTPContext(key.AddKey), middlewares.ValidateSignatureMessageMiddleware)
	r.POST("/did/:did/keys/reset", core.WithHTTPContext(key.ResetKey), middlewares.ValidateSignatureMessageMiddleware)
	r.POST("/did/:did/keys/:key_id/revoke", core.WithHTTPContext(key.RevokeKey), middlewares.ValidateSignatureMessageMiddleware)
}
