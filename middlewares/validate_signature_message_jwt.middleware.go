package middlewares

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/helpers"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/requests"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
)

type jwtPayload struct {
	core.BaseValidator
	Message *string `json:"message"`
}

func (r jwtPayload) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Message, "message"))

	return r.Error()
}

func ValidateSignatureJWTMessageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)

		payloadData := &jwtPayload{}
		if err := cc.BindWithValidate(payloadData); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}

		tokenM, _ := helpers.JWTVCDecodingT(utils.GetString(payloadData.Message), []byte(""))
		if tokenM == nil || tokenM.Header == nil || tokenM.Claims == nil || tokenM.Signature == "" {
			return c.JSON(emsgs.JWTInValid.GetStatus(), emsgs.JWTInValid.JSON())
		}

		msgPayload := &requests.VCJWTMessage{}
		_ = utils.MapToStruct(tokenM, &msgPayload)
		if err := cc.BindWithValidate(msgPayload); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}

		var isSigValid = false
		msgs := strings.Split(utils.GetString(payloadData.Message), ".")
		keyService := services.NewKeyService(cc)
		_, ierr := keyService.FindVerifiedHistoryPublicKey(msgPayload.Claims.Iss, fmt.Sprintf("%s.%s", msgs[0], msgs[1]), utils.GetString(msgPayload.Signature))
		if ierr == nil {
			isSigValid = true
		}

		if !isSigValid {
			return c.JSON(emsgs.SignatureInValid.GetStatus(), emsgs.SignatureInValid.JSON())
		}
		c.Set(consts.ContextKeyJWTData, msgPayload)
		c.Set(consts.ContextKeyMessage, utils.GetString(payloadData.Message))
		c.Set(consts.ContextKeySignature, utils.GetString(msgPayload.Signature))

		return next(c)
	}
}
