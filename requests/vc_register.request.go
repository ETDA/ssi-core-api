package requests

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type RequestVCRegister struct {
	core.BaseValidator
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
	Nonce      *string `json:"nonce"`
}

func (r RequestVCRegister) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Operation, "operation")) {
		r.Must(r.IsStrIn(r.Operation, consts.OperationVCRegister, "operation"))
	}

	r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.DIDAddress)}, "did_address"))
	r.Must(r.IsStrRequired(r.DIDAddress, "did_address"))


	if r.Must(r.IsStrRequired(r.Nonce, "nonce")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{
			"address": utils.GetString(r.DIDAddress),
			"nonce":   utils.GetString(r.Nonce),
		}, "nonce"))
	}
	return r.Error()
}
