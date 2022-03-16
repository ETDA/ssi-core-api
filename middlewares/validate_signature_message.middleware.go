package middlewares

import (
	"crypto/x509"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/requests"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
	"ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type payload struct {
	core.BaseValidator
	Message *string `json:"message"`
}

func (r payload) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Message, "message")) {
		r.Must(r.IsBase64(r.Message, "message"))
	}

	return r.Error()
}

func ValidateSignatureMessageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)

		payloadData := &payload{}
		if cc.GetSignature() == "" {
			return c.JSON(http.StatusBadRequest, core.NewValidatorFields(core.RequiredM("x-signature")))
		}
		if err := cc.BindWithValidate(payloadData); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}

		c.Set("message", utils.GetString(payloadData.Message))

		msgPayload := &requests.TxMessagePayload{}
		if err := cc.BindWithValidateMessage(msgPayload); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}
		operation := utils.GetString(msgPayload.Operation)
		keyType := utils.GetString(msgPayload.KeyType)

		var isSigValid = false
		if operation == consts.OperationDIDRegister {
			var algorithm x509.SignatureAlgorithm
			if keyType == consts.KeyTypeRSA2018 {
				algorithm = x509.SHA256WithRSA
			} else {
				algorithm = x509.ECDSAWithSHA256
			}

			isSigValid, _ = utils.VerifySignatureWithOption(
				utils.GetString(msgPayload.PublicKey),
				cc.GetSignature(),
				cc.GetMessage(),
				&utils.VerifySignatureOption{
					Algorithm: algorithm,
				})
		} else {
			keyService := services.NewKeyService(cc)
			_, ierr := keyService.FindVerifiedPublicKey(utils.GetString(msgPayload.DIDAddress), utils.GetString(payloadData.Message), cc.GetSignature())
			if ierr == nil {
				isSigValid = true
			}
		}
		if !isSigValid {
			return c.JSON(emsgs.SignatureInValid.GetStatus(), emsgs.SignatureInValid.JSON())
		}

		return next(c)
	}
}

func ValidateSignatureMessage(ctx core.IContext, tx []byte) error {
	request := &requests.TXBroadcastPayload{}
	err := utils.JSONParse(tx, request)
	if err != nil {
		return err
	}

	if err = request.Valid(ctx); err != nil {
		return err
	}

	messageBase64, err := utils.Base64Decode(utils.GetString(request.Message))
	if err != nil {
		return err
	}

	msgPayload := &requests.TxMessagePayload{}
	err = utils.JSONParse(utils.StringToBytes(messageBase64), msgPayload)
	if err != nil {
		return err
	}

	if err = msgPayload.Valid(ctx); err != nil {
		return err
	}

	var isSigValid = false
	if utils.GetString(msgPayload.Operation) == consts.OperationDIDRegister {
		var algorithm x509.SignatureAlgorithm
		if utils.GetString(msgPayload.KeyType) == consts.KeyTypeRSA2018 {
			algorithm = x509.SHA256WithRSA
		} else {
			algorithm = x509.ECDSAWithSHA256
		}

		isSigValid, _ = utils.VerifySignatureWithOption(
			utils.GetString(msgPayload.PublicKey),
			utils.GetString(request.Signature),
			utils.GetString(request.Message),
			&utils.VerifySignatureOption{
				Algorithm: algorithm,
			},
		)
	} else {
		keyService := services.NewKeyService(ctx)
		_, ierr := keyService.FindVerifiedPublicKey(utils.GetString(msgPayload.DIDAddress), utils.GetString(request.Message), utils.GetString(request.Signature))
		if ierr == nil {
			isSigValid = true
		}
	}

	if !isSigValid {
		return errors.New("INVALID_SIGNATURE")
	}

	return nil
}
