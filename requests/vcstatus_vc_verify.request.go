package requests

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"time"
)

type VCVerify struct {
	core.BaseValidator
	Context           []string   `json:"@context"`
	CredentialSubject core.Map   `json:"credentialSubject"`
	IssuanceDate      *time.Time `json:"issuanceDate"`
	Type              []string   `json:"type"`
	ID                *string    `json:"id"`
	Issuer            *string    `json:"issuer"`
}

func (r VCVerify) Valid(ctx core.IContext) core.IError {
	return r.Error()
}
