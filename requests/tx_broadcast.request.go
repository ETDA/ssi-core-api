package requests

import (
	"ssi-gitlab.teda.th/ssi/core"
	"time"
)

type TXBroadcastPayload struct {
	core.BaseValidator
	Message   *string    `json:"message"`
	Signature *string    `json:"signature"`
	Version   *string    `json:"version"`
	CreatedAt *time.Time `json:"created_at"`
}

func (r TXBroadcastPayload) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Message, "message"))
	r.Must(r.IsStrRequired(r.Signature, "signature"))
	r.Must(r.IsStrRequired(r.Version, "version"))
	r.Must(r.IsTimeRequired(r.CreatedAt, "created_at"))
	return r.Error()
}
