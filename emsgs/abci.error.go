package emsgs

import (
	"ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var (
	BroadcastTXError = core.Error{
		Status:  http.StatusBadGateway,
		Code:    "BROADCAST_TX_ERROR",
		Message: "Broadcast transaction error"}

	BroadcastTXErrorWithMSG = func(msg string) core.IError {
		return core.Error{
			Status:  http.StatusBadGateway,
			Code:    "BROADCAST_TX_ERROR",
			Message: msg}
	}
)
