package views

import (
	"time"
)

type VPVerify struct {
	VerificationResult bool       `json:"verification_result"`
	ID                 string     `json:"id"`
	Audience           string     `json:"audience"`
	Issuer             string     `json:"issuer"`
	IssuanceDate       *time.Time `json:"issuance_date"`
	ExpireDate         *time.Time `json:"expiration_date"`
	VC                 []VCVerify `json:"vc"`
}

func NewVPVerify(valid bool, id string, vc []VCVerify, issuer string, audience string, issuanceDateUnix int64, expirationDateUnix int64) *VPVerify {
	issuanceDate := time.Unix(issuanceDateUnix, 0)
	var expirationDate *time.Time = nil
	if expirationDateUnix != 0 {
		dateTime := time.Unix(expirationDateUnix, 0)
		expirationDate = &dateTime
	}
	return &VPVerify{
		VerificationResult: valid,
		ID:                 id,
		Audience:           audience,
		Issuer:             issuer,
		IssuanceDate:       &issuanceDate,
		ExpireDate:         expirationDate,
		VC:                 vc,
	}
}
