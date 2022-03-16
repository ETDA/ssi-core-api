package views

import (
	"time"

	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type VCVerify struct {
	VerificationResult       bool                   `json:"verification_result"`
	CID                      string                 `json:"cid"`
	Status                   *string                `json:"status"`
	IssuanceDate             *time.Time             `json:"issuance_date"`
	RevokeDate               *time.Time             `json:"revoke_date,omitempty"`
	ExpireDate               *time.Time             `json:"expire_date,omitempty"`
	Type                     []string               `json:"type"`
	Tags                     []string               `json:"tags"`
	Issuer                   string                 `json:"issuer"`
	Holder                   string                 `json:"holder"`
	SchemaVerificationResult map[string]interface{} `json:"schema_validation_result,omitempty"`
}

func NewVCVerify(vc *models.VC, valid bool, vcTypes []string, issuer string, holder string, issuanceDateUnix int64, expirationDateUnix int64, schemaValidateMessage interface{}) *VCVerify {
	issuanceDate := time.Unix(issuanceDateUnix, 0)
	var expirationDate *time.Time = nil
	if expirationDateUnix != 0 {
		dateTime := time.Unix(expirationDateUnix, 0)
		expirationDate = &dateTime
	}
	view := &VCVerify{
		VerificationResult: valid,
		CID:                vc.CID,
		Status:             vc.Status,
		IssuanceDate:       &issuanceDate,
		RevokeDate:         vc.RevokedAt,
		ExpireDate:         expirationDate,
		Type:               vcTypes,
		Tags:               vc.Tags,
		Issuer:             issuer,
		Holder:             holder,
	}
	if schemaValidateMessage != nil {
		var mapper map[string]interface{}
		_ = utils.JSONParse([]byte(utils.JSONToString(schemaValidateMessage)), &mapper)
		view.SchemaVerificationResult = mapper
		return view
	}

	return view
}
