package requests

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type VCStatusAdd struct {
	core.BaseValidator
	Status     *string `json:"status"`
	CID        *string `json:"cid"`
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
	VCHash     *string `json:"vc_hash"`
	Nonce      *string `json:"nonce"`
}

func (r VCStatusAdd) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Status, "status"))
	r.Must(r.IsStrIn(r.Status, consts.VCStatusActive, "status"))
	if r.Must(r.IsStrRequired(r.CID, "cid")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.VC{}).TableName(), bson.M{"cid": utils.GetString(r.CID)}, "cid"))
	}

	r.Must(r.IsStrRequired(r.Operation, "operation"))
	r.Must(r.IsStrRequired(r.VCHash, "vc_hash"))
	r.Must(r.IsStrIn(r.Operation, consts.OperationVCAddStatus, "operation"))
	r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": *r.DIDAddress}, "did_address"))

	if r.Must(r.IsStrRequired(r.Nonce, "nonce")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{
			"address": utils.GetString(r.DIDAddress),
			"nonce":   utils.GetString(r.Nonce),
		}, "nonce"))
	}
	return r.Error()
}
