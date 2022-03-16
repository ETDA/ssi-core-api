package services

import (
	"crypto/x509"
	"errors"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type IKeyService interface {
	Add(payload *KeyAdd) (*models.DID, core.IError)
	FindByPublicKey(didAddress string, nextKeyHash string) (*models.Key, core.IError)
	Revoke(payload *KeyRevoke) (*models.DID, core.IError)
	ResetByRecoverer(payload *KeyReset) (*models.DID, core.IError)
	Find(id string) (*models.Key, core.IError)
	FindVerifiedPublicKey(didAddress string, message string, signature string) (*models.Key, core.IError)
	FindVerifiedHistoryPublicKey(didAddress string, message string, signature string) (*models.Key, core.IError)
}
type KeyService struct {
	ctx        core.IContext
	didService IDIDService
}

func NewKeyService(ctx core.IContext) IKeyService {
	return &KeyService{ctx, NewDIDService(ctx)}
}

type KeyAdd struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type KeyReset struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (s KeyService) Add(payload *KeyAdd) (*models.DID, core.IError) {
	abciSvc := NewABCIService(s.ctx)
	res, ierr := abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	didResult := &models.DID{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), didResult)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return didResult, nil
}

func (s KeyService) FindByPublicKey(didAddress string, publicKey string) (*models.Key, core.IError) {
	key := &models.Key{}

	if err := s.ctx.DBMongo().FindAggregateOne(key, models.DID{}.TableName(), []bson.M{
		{"$unwind": "$keys"},
		{
			"$match": bson.M{
				"address":                     didAddress,
				"keys.public_keys.used_at":    nil,
				"keys.public_keys.public_key": publicKey,
			}},
		{
			"$replaceRoot": bson.M{"newRoot": "$keys"},
		},
	}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, s.ctx.NewError(err, emsgs.NotFound)
		}

		return nil, s.ctx.NewError(err, emsgs.DBError, didAddress, publicKey)
	}

	return key, nil
}

type KeyRevoke struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (s KeyService) Revoke(payload *KeyRevoke) (*models.DID, core.IError) {
	abciSvc := NewABCIService(s.ctx)

	res, ierr := abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	didResult := &models.DID{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), didResult)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return s.didService.Find(didResult.Address)
}

func (s KeyService) ResetByRecoverer(payload *KeyReset) (*models.DID, core.IError) {
	abciSvc := NewABCIService(s.ctx)

	res, ierr := abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	didResult := &models.DID{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), didResult)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return s.didService.Find(didResult.Address)
}

func (s KeyService) Find(id string) (*models.Key, core.IError) {
	key := &models.Key{}

	if err := s.ctx.DBMongo().FindAggregateOne(key, models.DID{}.TableName(), []bson.M{{
		"$match": bson.M{
			"keys._id": id,
		}},
		{
			"$project": bson.M{
				"_id": 0,
				"keys": bson.M{
					"$filter": bson.M{
						"input": "$keys",
						"as":    "key",
						"cond":  bson.M{"$eq": []string{"$$key._id", id}},
					},
				},
			},
		},
		{"$unwind": "$keys"},
		{"$replaceRoot": bson.M{"newRoot": "$keys"}},
	}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, s.ctx.NewError(err, emsgs.NotFound)
		}

		return nil, s.ctx.NewError(err, emsgs.DBError, id)
	}

	return key, nil

}
func (s KeyService) FindVerifiedPublicKey(didAddress string, message string, signature string) (*models.Key, core.IError) {
	did, ierr := s.didService.Find(didAddress)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	for _, key := range did.Keys {
		if key.RevokedAt != nil {
			continue
		}
		var algorithm x509.SignatureAlgorithm
		if key.KeyType == consts.KeyTypeRSA2018 {
			algorithm = x509.SHA256WithRSA
		} else {
			algorithm = x509.ECDSAWithSHA256
		}

		isTempSigValid, _ := utils.VerifySignatureWithOption(
			key.PublicKey,
			signature,
			message,
			&utils.VerifySignatureOption{
				Algorithm: algorithm,
			})

		if isTempSigValid {
			return &key, nil
		}
	}
	return nil, s.ctx.NewError(emsgs.NotFound, emsgs.NotFound)
}

func (s KeyService) FindVerifiedHistoryPublicKey(didAddress string, message string, signature string) (*models.Key, core.IError) {
	did, ierr := s.didService.Find(didAddress)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	for _, didHistory := range did.History {
		for _, key := range didHistory.Keys {
			if key.RevokedAt != nil {
				continue
			}
			var algorithm x509.SignatureAlgorithm
			if key.KeyType == consts.KeyTypeRSA2018 {
				algorithm = x509.SHA256WithRSA
			} else {
				algorithm = x509.ECDSAWithSHA256
			}

			isTempSigValid, _ := utils.VerifySignatureWithOption(
				key.PublicKey,
				signature,
				message,
				&utils.VerifySignatureOption{
					Algorithm: algorithm,
				})

			if isTempSigValid {
				return &key, nil
			}
		}
	}
	return nil, s.ctx.NewError(emsgs.NotFound, emsgs.NotFound)
}
