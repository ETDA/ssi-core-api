package requests

import (
	"fmt"
	"strings"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type TxMessagePayload struct {
	core.BaseValidator
	PublicKey  *string `json:"public_key"`
	Nonce      *string `json:"nonce"`
	Operation  *string `json:"operation"`
	DIDAddress *string `json:"did_address"`
	KeyType    *string `json:"key_type"`
}

func (r TxMessagePayload) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Operation, "operation")) &&
		r.Must(r.IsStrIn(r.Operation, fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
			consts.OperationDIDRegister,
			consts.OperationDIDKeyAdd,
			consts.OperationDIDKeyReset,
			consts.OperationDIDRecovererAdd,
			consts.OperationDIDKeyRevoke,
			consts.OperationVCRegister,
			consts.OperationVCAddStatus,
			consts.OperationVCUpdateStatus,
			consts.OperationVCTagStatus,
		), "operation")) {

		didTableName := (&models.DID{}).TableName()
		if utils.GetString(r.Operation) == consts.OperationDIDRegister {
			r.Must(r.IsStrRequired(r.PublicKey, "public_key"))
			r.Must(r.IsStrRequired(r.KeyType, "key_type"))
			r.Must(r.IsStrIn(r.KeyType, strings.Join(consts.SupportedKeys, "|"), "key_type"))
		}
		if utils.GetString(r.Operation) != consts.OperationDIDRegister {
			if r.Must(r.IsStrRequired(r.DIDAddress, "did_address")) &&
				r.Must(r.IsStrRequired(r.Nonce, "nonce")) {

				if r.Must(r.IsMongoExistsWithCondition(ctx, didTableName, bson.M{"address": utils.GetString(r.DIDAddress)}, "did_address")) &&
					r.Must(r.IsMongoExistsWithCondition(ctx, didTableName, bson.M{"nonce": utils.GetString(r.Nonce), "address": utils.GetString(r.DIDAddress)}, "nonce")) {
				}
			}
		}
	}
	return r.Error()
}
