package requests

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type RequestDIDRecovererAdd struct {
	core.BaseValidator
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
	Recoverer  *string `json:"recoverer"`
	Nonce      *string `json:"nonce"`
}

func (r *RequestDIDRecovererAdd) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Operation, "operation"))
	r.Must(r.IsStrIn(r.Operation, consts.OperationDIDRecovererAdd, "operation"))

	if r.Must(r.IsStrRequired(r.DIDAddress, "did_address")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.DIDAddress)}, "did_address"))
	}

	if r.Must(r.IsStrRequired(r.Recoverer, "recoverer")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.Recoverer)}, "recoverer"))
	}

	if r.Must(r.IsStrRequired(r.Nonce, "nonce")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{
			"address": utils.GetString(r.DIDAddress),
			"nonce":   utils.GetString(r.Nonce),
		}, "nonce"))
	}

	return r.Error()
}
