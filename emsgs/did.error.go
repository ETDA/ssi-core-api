package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
)

var (
	RecovererIsNotMatchError = &core.IValidMessage{
		Name:    "request_did",
		Code:    "RECOVERER_IS_NOT_MATCH",
		Message: "your did address is not match with request's recoverer did address",
	}
)
