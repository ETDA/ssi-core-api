package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var (
	VCCannotUpdateStatusBeforeAdd = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "VC_CAN_NOT_UPDATE",
		Message: "can't update. please add this VC status first"}
	VCCannotTagStatusNotActive = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "VC_CAN_NOT_TAG",
		Message: "Can't Tag status. Please provide VC which status active only."}
	VCCannotTagStatusMultiple = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "VC_CAN_NOT_TAG",
		Message: "Can't Tag status. This VC was tagged previously."}
)
