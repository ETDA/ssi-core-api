package did

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/requests"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/views"
	"ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

type VCController struct{}

func (n *VCController) Register(c core.IHTTPContext) error {
	input := &requests.RequestVCRegister{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	service := services.NewVCService(c)
	vc, err := service.Register(&services.VCRegister{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	})
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusCreated, views.NewVC(vc))
}

func (n *VCController) Find(c core.IHTTPContext) error {
	service := services.NewVCService(c)
	vc, err := service.Find(c.Param("id"))
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	return c.JSON(http.StatusOK, views.NewVC(vc))
}
