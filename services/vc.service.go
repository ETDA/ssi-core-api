package services

import (
	"errors"
	"time"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type IVCService interface {
	Register(payload *VCRegister) (*models.VC, core.IError)
	Update(payload *VCRevoke) (*models.VC, core.IError)
	Find(cid string) (*models.VC, core.IError)
	FindMultiple(cids []string) ([]models.VC, core.IError)
	AddStatus(payload *VCStatusAddPayload) (*models.VC, core.IError)
	UpdateStatus(id string, payload *VCStatusUpdatePayload) (*models.VC, core.IError)
	TagStatus(payload *VCStatusTagPayload) ([]models.VC, core.IError)
}

type VCService struct {
	ctx        core.IContext
	didService IDIDService
	abciSvc    IABCIService
}

func NewVCService(ctx core.IContext) IVCService {
	return &VCService{ctx, NewDIDService(ctx), NewABCIService(ctx)}
}

type VCRegister struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (s VCService) Find(cid string) (*models.VC, core.IError) {
	vc := &models.VC{}
	err := s.ctx.DBMongo().FindOne(vc, vc.TableName(), bson.M{
		"cid": cid,
	})

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, s.ctx.NewError(err, emsgs.NotFound)
	}

	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError, cid)
	}

	return vc, nil
}

func (s VCService) FindMultiple(cids []string) ([]models.VC, core.IError) {
	vc := make([]models.VC, 0)
	bsonCIDs := bson.A{}
	for _, cid := range cids {
		bsonCIDs = append(bsonCIDs, cid)
	}
	err := s.ctx.DBMongo().Find(&vc, models.VC{}.TableName(), bson.M{
		"cid": bson.M{
			"$in": bsonCIDs,
		},
	})

	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError, bsonCIDs)
	}

	return vc, nil
}

func (s VCService) Register(payload *VCRegister) (*models.VC, core.IError) {
	res, ierr := s.abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	vcResult := &models.VC{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), vcResult)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return s.Find(vcResult.CID)
}

type VCRevoke struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func (s VCService) Update(payload *VCRevoke) (*models.VC, core.IError) {
	res, ierr := s.abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	vcResult := &models.VC{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), vcResult)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return s.Find(vcResult.CID)
}

func (s VCService) UpdateStatus(id string, payload *VCStatusUpdatePayload) (*models.VC, core.IError) {
	vc, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr, id)
	}

	if vc.DIDAddress != payload.DIDAddress {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound, id)
	}

	if vc.Status == nil {
		return nil, s.ctx.NewError(emsgs.VCCannotUpdateStatusBeforeAdd, emsgs.VCCannotUpdateStatusBeforeAdd, id)
	}

	res, ierr := s.abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	status := &models.VC{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), status)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return s.Find(id)
}

func (s VCService) AddStatus(payload *VCStatusAddPayload) (*models.VC, core.IError) {
	vc, ierr := s.Find(payload.CID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr, payload.CID)
	}

	if vc.DIDAddress != payload.DIDAddress {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound, payload.CID)
	}

	res, ierr := s.abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	status := &models.VC{}
	err := utils.JSONParse(utils.StringToBytes(res.Result.DeliverTx.Info), status)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.InternalServerError, res)
	}

	return s.Find(status.CID)
}

func (s VCService) TagStatus(payload *VCStatusTagPayload) ([]models.VC, core.IError) {

	for _, cid := range payload.CIDs {
		vc, ierr := s.Find(cid)
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr, payload.CIDs)
		}
		if utils.GetString(vc.Status) != consts.VCStatusActive {
			return nil, s.ctx.NewError(emsgs.VCCannotTagStatusNotActive, emsgs.VCCannotTagStatusNotActive, cid)

		}
	}
	_, ierr := s.abciSvc.BroadcastTXCommit(&TXBroadcastPayload{
		Message:   payload.Message,
		Signature: payload.Signature,
		Version:   consts.BlockChainVersion100,
		CreatedAt: utils.GetCurrentDateTime(),
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	status := make([]models.VC, 0)
	for _, cid := range payload.CIDs {
		vc, _ := s.Find(cid)
		status = append(status, *vc)
	}

	return status, nil
}

type VCStatusVerifyProofPayload struct {
	Signature          string     `json:"signature"`
	ProofPurpose       string     `json:"proofPurpose"`
	Type               string     `json:"type"`
	VerificationMethod string     `json:"verificationMethod"`
	Created            *time.Time `json:"created"`
}

type VCStatusVerifyPayload struct {
	CID     string
	Issuer  string
	Message string
}

type VCStatusUpdatePayload struct {
	Status     string `json:"status"`
	CID        string `json:"cid"`
	DIDAddress string `json:"did_address"`
	Message    string `json:"message"`
	Signature  string `json:"signature"`
}

type VCStatusAddPayload struct {
	Status     string `json:"status"`
	CID        string `json:"cid"`
	DIDAddress string `json:"did_address"`
	Message    string `json:"message"`
	Signature  string `json:"signature"`
}

type VCStatusTagPayload struct {
	CIDs       []string `json:"cids"`
	Tags       []string `json:"tags"`
	DIDAddress string   `json:"did_address"`
	Message    string   `json:"message"`
	Signature  string   `json:"signature"`
}
