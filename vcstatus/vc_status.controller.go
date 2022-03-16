package vcstatus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/helpers"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/requests"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/views"
)

type VCStatusController struct{}

func (n *VCStatusController) All(c core.IHTTPContext) error {
	service := services.NewVCService(c)
	cid := c.QueryParam("cid")
	split := strings.Split(cid, ",")

	res, ierr := service.FindMultiple(split)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, res)
}

func (n *VCStatusController) Find(c core.IHTTPContext) error {
	service := services.NewVCService(c)
	id := c.Param("id")

	res, ierr := service.Find(id)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, views.NewVCStatus(res))
}

func (n *VCStatusController) Update(c core.IHTTPContext) error {
	input := &requests.VCStatusUpdate{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	if utils.GetString(input.CID) != c.Param("id") {
		return c.JSON(emsgs.BadRequest.GetStatus(), emsgs.BadRequest.JSON())
	}
	service := services.NewVCService(c)
	payload := &services.VCStatusUpdatePayload{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	}
	_ = utils.Copy(payload, input)
	vc, ierr := service.UpdateStatus(c.Param("id"), payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, views.NewVCStatus(vc))
}

func (n *VCStatusController) Add(c core.IHTTPContext) error {
	input := &requests.VCStatusAdd{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	service := services.NewVCService(c)
	payload := &services.VCStatusAddPayload{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	}
	_ = utils.Copy(payload, input)
	vc, ierr := service.AddStatus(payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, views.NewVCStatus(vc))
}

func (n *VCStatusController) VerifyVC(c core.IHTTPContext) error {
	jwt := c.Get(consts.ContextKeyJWTData).(*requests.VCJWTMessage)
	service := services.NewVCService(c)
	res, ierr := service.Find(jwt.Claims.Jti)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	vcTypes := jwt.Claims.VC.Type
	credentialSubject := json.RawMessage(utils.JSONToString(jwt.Claims.VC.CredentialSubject))
	credentialSchemaID := utils.GetString(jwt.Claims.VC.CredentialSchema.ID)
	schemaService := services.NewVCSchemaService(c)
	valid, message, ierr := schemaService.Validate(&services.ValidateVCSchemaSchemaPayload{
		ID:       credentialSchemaID,
		Document: &credentialSubject,
		VCTypes:  vcTypes,
	})
	if ierr != nil {
		return c.JSON(http.StatusOK, views.NewVCVerify(res, valid, vcTypes, jwt.Claims.Iss, jwt.Claims.Sub, jwt.Claims.Nbf, jwt.Claims.Exp, ierr.JSON()))
	}

	return c.JSON(http.StatusOK, views.NewVCVerify(res, valid, vcTypes, jwt.Claims.Iss, jwt.Claims.Sub, jwt.Claims.Nbf, jwt.Claims.Exp, message))

}

func (n *VCStatusController) Tag(c core.IHTTPContext) error {
	input := &requests.VCStatusTag{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	service := services.NewVCService(c)
	// payload := &services.VCStatus
	payload := &services.VCStatusTagPayload{
		Message:   c.GetMessage(),
		Signature: c.GetSignature(),
	}
	_ = utils.Copy(payload, input)
	vcs, ierr := service.TagStatus(payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.Map{"vcs": views.NewVCStatusList(vcs)})
}

func (n *VCStatusController) VerifyVP(c core.IHTTPContext) error {
	jwt := c.Get(consts.ContextKeyJWTData).(*requests.VCJWTMessage)
	service := services.NewVCService(c)

	vcList := make([]views.VCVerify, 0)
	for _, vcJWT := range jwt.Claims.VP.VerifiableCredential {
		tokenM, _ := helpers.JWTVCDecodingT(vcJWT, []byte(""))
		if tokenM == nil || tokenM.Header == nil || tokenM.Claims == nil || tokenM.Signature == "" {
			return c.JSON(emsgs.JWTInValid.GetStatus(), emsgs.JWTInValid.JSON())
		}

		msgPayload := &requests.VCJWTMessage{}
		_ = utils.MapToStruct(tokenM, msgPayload)
		if err := c.BindWithValidate(msgPayload); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}

		isSigValid := false
		msgs := strings.Split(vcJWT, ".")
		keyService := services.NewKeyService(c)
		_, ierr := keyService.FindVerifiedHistoryPublicKey(msgPayload.Claims.Iss, fmt.Sprintf("%s.%s", msgs[0], msgs[1]), utils.GetString(msgPayload.Signature))
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		isSigValid = true
		var vc *models.VC
		if isSigValid {
			res, ierr := service.Find(msgPayload.Claims.Jti)
			if ierr != nil {
				return c.JSON(ierr.GetStatus(), ierr.JSON())
			}

			vc = res
		} else {
			expiredAt := time.Unix(msgPayload.Claims.Exp, 0).UTC()
			createdAt := time.Unix(msgPayload.Claims.Iat, 0).UTC()

			vc = &models.VC{
				CID:        msgPayload.Claims.Jti,
				DIDAddress: msgPayload.Claims.Iss,
				CreatedAt:  &createdAt,
				ExpiredAt:  &expiredAt,
			}
		}

		vcTypes := msgPayload.Claims.VC.Type
		credentialSubject := json.RawMessage(utils.JSONToString(msgPayload.Claims.VC.CredentialSubject))
		credentialSchemaID := utils.GetString(msgPayload.Claims.VC.CredentialSchema.ID)

		schemaService := services.NewVCSchemaService(c)
		isSchemaValid, message, ierr := schemaService.Validate(&services.ValidateVCSchemaSchemaPayload{
			ID:       credentialSchemaID,
			Document: &credentialSubject,
			VCTypes:  vcTypes,
		})
		if ierr != nil {
			vcList = append(vcList, *views.NewVCVerify(vc, isSigValid && isSchemaValid, msgPayload.Claims.VC.Type, msgPayload.Claims.Iss, msgPayload.Claims.Sub, msgPayload.Claims.Nbf, msgPayload.Claims.Exp, ierr.JSON()))
			continue
		}
		vcList = append(vcList, *views.NewVCVerify(vc, isSigValid && isSchemaValid, msgPayload.Claims.VC.Type, msgPayload.Claims.Iss, msgPayload.Claims.Sub, msgPayload.Claims.Nbf, msgPayload.Claims.Exp, message))

	}

	return c.JSON(http.StatusOK, views.NewVPVerify(true, jwt.Claims.Jti, vcList, jwt.Claims.Iss, jwt.Claims.Aud, jwt.Claims.Nbf, jwt.Claims.Exp))
}
