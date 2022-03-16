package requests

import (
	"strings"

	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type RequestDIDRegister struct {
	core.BaseValidator
	Operation *string `json:"operation"`
	PublicKey *string `json:"public_key"`
	KeyType   *string `json:"key_type"`
}

func (r RequestDIDRegister) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Operation, "operation")) {
		r.Must(r.IsStrIn(r.Operation, consts.OperationDIDRegister, "operation"))
	}
	r.Must(r.IsStrRequired(r.KeyType, "key_type"))
	r.Must(r.IsStrIn(r.KeyType, strings.Join(consts.SupportedKeys, "|"), "key_type"))

	r.Must(r.IsMongoStrUnique(ctx, (&models.DID{}).TableName(), bson.M{
		"address": utils.GenerateDID(utils.NewSha256(*r.PublicKey), ctx.ENV().Config().DIDMethodDefault),
	}, "did_address"))

	r.Must(r.IsStrRequired(r.PublicKey, "public_key"))
	return r.Error()
}
