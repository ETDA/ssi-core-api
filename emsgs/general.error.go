package emsgs

import (
	"ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var (
	DBError = core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "DATABASE_ERROR",
		Message: "database internal error"}

	InternalServerError = core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error"}

	NotFound = core.Error{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "not found"}

	BadRequest = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: "bad request"}

	SignatureInValid = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_SIGNATURE",
		Message: "Signature is not valid"}

	JWTInValid = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_JWT",
		Message: "JWT is not valid"}

	JSONInValid = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_JSON",
		Message: "Must be json format"}

	Unauthorized = core.Error{
		Status:  http.StatusUnauthorized,
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized"}
)
