package did

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/requests"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/views"
)

type KeyController struct{}

func (n *KeyController) AddKey(c core.IHTTPContext) error {
	input := &requests.RequestDIDKeyAdd{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	if utils.GetString(input.DIDAddress) != c.Param("did") {
		return c.JSON(emsgs.BadRequest.GetStatus(), emsgs.BadRequest.JSON())
	}
	service := services.NewKeyService(c)
	res, err := service.Add(&services.KeyAdd{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	})

	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusCreated, views.NewDIDDocument(res))
}

func (n *KeyController) RevokeKey(c core.IHTTPContext) error {
	input := &requests.RequestDIDRevokeKey{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	if utils.GetString(input.DIDAddress) != c.Param("did") {
		return c.JSON(emsgs.BadRequest.GetStatus(), emsgs.BadRequest.JSON())
	}
	if utils.GetString(input.KeyID) != c.Param("key_id") {
		return c.JSON(emsgs.BadRequest.GetStatus(), emsgs.BadRequest.JSON())
	}
	service := services.NewKeyService(c)
	res, err := service.Revoke(&services.KeyRevoke{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	})
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusOK, views.NewDIDDocument(res))
}

func (n *KeyController) ResetKey(c core.IHTTPContext) error {
	input := &requests.RequestDIDKeyReset{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	if utils.GetString(input.DIDAddress) != c.Param("did") {
		return c.JSON(emsgs.BadRequest.GetStatus(), emsgs.BadRequest.JSON())
	}
	service := services.NewKeyService(c)
	res, err := service.ResetByRecoverer(&services.KeyReset{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	})
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusOK, views.NewDIDDocument(res))
}
