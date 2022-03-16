package services

import (
	"errors"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDIDService interface {
	Find(address string) (*models.DID, core.IError)
	Register(payload *DIDRegister) (*models.DID, core.IError)
	AddRecoverer(payload *AddRecovererPayload) (*models.DID, core.IError)
	GetNonce(didAddress string) (string, core.IError)
	FindByVersion(address string, version string) (*models.DIDHistory, core.IError)
}

type DIDService struct {
	ctx core.IContext
}

func NewDIDService(ctx core.IContext) IDIDService {
	return &DIDService{ctx}
}

func (s DIDService) Find(address string) (*models.DID, core.IError) {
	did := &models.DID{}

	if err := s.ctx.DBMongo().FindOne(did, did.TableName(), bson.M{
		"address": address,
	}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, s.ctx.NewError(err, emsgs.NotFound)
		}

		return nil, s.ctx.NewError(err, emsgs.DBError, address)
	}
	return did, nil
}

type DIDRegister struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type AddRecovererPayload struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (s DIDService) Register(payload *DIDRegister) (*models.DID, core.IError) {
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

	return s.Find(didResult.Address)
}

func (s DIDService) AddRecoverer(payload *AddRecovererPayload) (*models.DID, core.IError) {
	// TODO: Check recoveror must not exists

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

func (s DIDService) GetNonce(didAddress string) (string, core.IError) {
	did := &models.DID{}
	if err := s.ctx.DBMongo().FindOne(did, did.TableName(), bson.M{"address": didAddress}, &options.FindOneOptions{
		Projection: bson.M{
			"nonce": 1,
		},
	}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", s.ctx.NewError(err, emsgs.NotFound)
		}

		return "", s.ctx.NewError(err, emsgs.DBError)
	}

	return did.Nonce, nil
}

func (s DIDService) FindByVersion(address string, version string) (*models.DIDHistory, core.IError) {
	did := &models.DID{}

	if err := s.ctx.DBMongo().FindOne(did, did.TableName(), bson.M{
		"address": address,
	}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, s.ctx.NewError(err, emsgs.NotFound)
		}

		return nil, s.ctx.NewError(err, emsgs.DBError, address)
	}
	for _, didHistory := range did.History {
		if didHistory.Version == version {
			return &didHistory, nil
		}
	}
	return nil, s.ctx.NewError(emsgs.NotFound, emsgs.NotFound)
}
