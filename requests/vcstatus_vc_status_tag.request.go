package requests

import (
	"fmt"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type VCStatusTag struct {
	core.BaseValidator
	CIDs       []string `json:"cids"`
	Tags       []string `json:"tags"`
	Operation  *string  `json:"operation"`
	DIDAddress *string  `json:"did_address"`
	Nonce      *string  `json:"nonce"`
}

func (r VCStatusTag) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsRequiredArray(r.CIDs, "cids")) {
		for index, cid := range r.CIDs {

			r.Must(r.IsMongoExistsWithCondition(ctx, (&models.VC{}).TableName(), bson.M{"cid": cid}, fmt.Sprintf("cid[%d]", index)))
		}
	}
	if r.Must(r.IsRequiredArray(r.Tags, "tags")) {
		for index, tag := range r.Tags {
			r.Must(r.IsStrRequired(&tag, fmt.Sprintf("tags[%d]", index)))
		}
	}

	r.Must(r.IsStrRequired(r.Operation, "operation"))
	r.Must(r.IsStrIn(r.Operation, consts.OperationVCTagStatus, "operation"))
	r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.DIDAddress)}, "did_address"))

	if r.Must(r.IsStrRequired(r.Nonce, "nonce")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{
			"address": utils.GetString(r.DIDAddress),
			"nonce":   utils.GetString(r.Nonce),
		}, "nonce"))
	}
	return r.Error()
}
