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

type DIDController struct{}

func (n *DIDController) Find(c core.IHTTPContext) error {
	didSvc := services.NewDIDService(c)
	did, err := didSvc.Find(c.Param("did"))
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusOK, views.NewDIDDocument(did))
}

func (n *DIDController) Register(c core.IHTTPContext) error {
	input := &requests.RequestDIDRegister{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	service := services.NewDIDService(c)
	did, err := service.Register(&services.DIDRegister{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	})
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusCreated, views.NewDIDDocument(did))
}

func (n *DIDController) AddRecoverer(c core.IHTTPContext) error {
	input := &requests.RequestDIDRecovererAdd{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	if utils.GetString(input.DIDAddress) != c.Param("did") {
		return c.JSON(emsgs.BadRequest.GetStatus(), emsgs.BadRequest.JSON())
	}
	service := services.NewDIDService(c)
	did, err := service.AddRecoverer(&services.AddRecovererPayload{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	})
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusCreated, views.NewDIDDocument(did))
}

func (n *DIDController) GetNonce(c core.IHTTPContext) error {
	service := services.NewDIDService(c)
	nonce, err := service.GetNonce(c.Param("did"))
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	return c.JSON(http.StatusOK, core.Map{
		"nonce": nonce,
	})
}

func (n *DIDController) History(c core.IHTTPContext) error {
	didSvc := services.NewDIDService(c)
	did, err := didSvc.Find(c.Param("did"))
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusOK, views.NewDIDDocumentHistory(did))
}
func (n *DIDController) FindByVersion(c core.IHTTPContext) error {
	didSvc := services.NewDIDService(c)
	did, err := didSvc.FindByVersion(c.Param("did"), c.Param("version_id"))
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusOK, views.NewDIDDocumentVersion(did))
}
