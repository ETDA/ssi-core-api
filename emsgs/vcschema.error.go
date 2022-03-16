package emsgs

import (
	"fmt"
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

func JSONSchemaNotFound(schemaID string) core.IError {
	return core.Error{
		Status:  http.StatusNotFound,
		Code:    "SCHEMA_NOT_FOUND",
		Message: fmt.Sprintf("the schema $id %v is not found or provided $id is not json schema.", schemaID),
	}
}
