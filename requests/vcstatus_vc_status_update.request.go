package requests

import (
	"fmt"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type VCStatusUpdate struct {
	core.BaseValidator
	Status     *string `json:"status"`
	CID        *string `json:"cid"`
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
	Nonce      *string `json:"nonce"`
}

func (r VCStatusUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Status, "status"))
	r.Must(r.IsStrIn(r.Status, fmt.Sprintf("%s|%s",
		consts.VCStatusExpired,
		consts.VCStatusRevoke,
	), "status"))
	if r.Must(r.IsStrRequired(r.CID, "cid")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.VC{}).TableName(), bson.M{"cid": utils.GetString(r.CID)}, "cid"))
	}

	r.Must(r.IsStrRequired(r.Operation, "operation"))
	r.Must(r.IsStrIn(r.Operation, consts.OperationVCUpdateStatus, "operation"))
	r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.DIDAddress)}, "did_address"))

	if r.Must(r.IsStrRequired(r.Nonce, "nonce")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{
			"address": utils.GetString(r.DIDAddress),
			"nonce":   utils.GetString(r.Nonce),
		}, "nonce"))
	}
	return r.Error()
}
