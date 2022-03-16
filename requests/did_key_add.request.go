package requests

import (
	"crypto/x509"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"strings"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type RequestDIDKeyAddNewKey struct {
	PublicKey  *string `json:"public_key"`
	Signature  *string `json:"signature"`
	Controller *string `json:"controller"`
	KeyType    *string `json:"key_type"`
}

type RequestDIDKeyAdd struct {
	core.BaseValidator
	Operation  *string                 `json:"operation"`
	DIDAddress *string                 `json:"did_address"`
	PublicKey  *string                 `json:"public_key"`
	NewKey     *RequestDIDKeyAddNewKey `json:"new_key"`
	Nonce      *string                 `json:"nonce"`
}

func (r RequestDIDKeyAdd) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Operation, "operation")) {
		r.Must(r.IsStrIn(r.Operation, consts.OperationDIDKeyAdd, "operation"))
	}
	if r.Must(r.IsStrRequired(r.DIDAddress, "did_address")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.DIDAddress)}, "did_address"))
	}
	if r.Must(r.IsRequired(r.NewKey, "new_key")) {
		r.Must(r.IsStrRequired(r.NewKey.PublicKey, "new_key.public_key"))
		r.Must(r.IsStrRequired(r.NewKey.Signature, "new_key.signature"))

		if r.NewKey.Controller != nil {
			r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{"address": utils.GetString(r.NewKey.Controller)}, "new_key.controller"))
		}

		if r.NewKey.PublicKey != nil && r.NewKey.Signature != nil {
			if r.Must(r.IsStrIn(r.NewKey.KeyType, strings.Join(consts.SupportedKeys, "|"), "new_key.key_type")) {
				var algorithm x509.SignatureAlgorithm
				if utils.GetString(r.NewKey.KeyType) == consts.KeyTypeRSA2018 {
					algorithm = x509.SHA256WithRSA
				} else {
					algorithm = x509.ECDSAWithSHA256
				}

				isSigValid, _ := utils.VerifySignatureWithOption(
					utils.GetString(r.NewKey.PublicKey),
					utils.GetString(r.NewKey.Signature),
					utils.GetString(r.NewKey.PublicKey),
					&utils.VerifySignatureOption{
						Algorithm: algorithm,
					},
				)

				if !isSigValid {
					r.Must(false, &core.IValidMessage{
						Name:    "new_key",
						Code:    emsgs.SignatureInValid.Code,
						Message: emsgs.SignatureInValid.Message.(string),
					})
				}
			}
		}
	}

	if r.Must(r.IsStrRequired(r.Nonce, "nonce")) {
		r.Must(r.IsMongoExistsWithCondition(ctx, (&models.DID{}).TableName(), bson.M{
			"address": utils.GetString(r.DIDAddress),
			"nonce":   utils.GetString(r.Nonce),
		}, "nonce"))
	}
	return r.Error()
}
