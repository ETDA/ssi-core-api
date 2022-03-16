package services

import (
	"time"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type BCTranasctionServiceOptions struct {
	BlockHashHex string
}

func NewBCTransactionService(ctx core.IContext, options *BCTranasctionServiceOptions) IBCTransactionService {
	return &BCTransactionService{ctx: ctx,
		options:    options,
		didService: NewDIDService(ctx),
		keyService: NewKeyService(ctx),
		vcService:  NewVCService(ctx),
	}
}

type IBCTransactionService interface {
	CreateBlock(hash string, height int64) (*models.Block, core.IError)
	ConfirmBlock(hash string, confirmedAt *time.Time) core.IError
	RegisterDID(data *models.DIDRegisterTX) (*models.DID, core.IError)
	DIDAddRecoverer(data *models.DIDRecovererAddTX) (*models.DID, core.IError)
	AddKey(data *models.DIDKeyAddTX) (*models.DID, core.IError)
	RevokeKey(data *models.DIDKeyRevokeTX) (*models.DID, core.IError)
	ResetKey(data *models.DIDKeyResetTX) (*models.DID, core.IError)
	VCRegister(data *models.VCRegisterTX) (*models.VC, core.IError)
	UpdateDIDDocumentVersion(signature string, did string, createdAt *time.Time) core.IError
	UpdateNonce(didAddress string, signature string) core.IError
	VCAddStatus(data *models.VCUpdateStatusTX) (*models.VC, core.IError)
	VCUpdateStatus(data *models.VCUpdateStatusTX) (*models.VC, core.IError)
	VCTagStatus(data *models.VCTagStatusTX) ([]models.VC, core.IError)
}

type BCTransactionService struct {
	ctx        core.IContext
	options    *BCTranasctionServiceOptions
	didService IDIDService
	keyService IKeyService
	vcService  IVCService
}

type createTransactionPayload struct {
	TX           []byte
	DIDAddress   string
	KeyID        string
	Operation    string
	Signature    string
	CreatedAt    *time.Time
	BlockHashHex *string
}

func (s BCTransactionService) createTransaction(data *createTransactionPayload) (*models.Transaction, core.IError) {
	transaction := &models.Transaction{
		ID:         utils.NewSha256(utils.BytesToString(data.TX)),
		DIDAddress: data.DIDAddress,
		KeyID:      data.KeyID,
		Message:    utils.BytesToString(data.TX),
		Operation:  data.Operation,
		Signature:  data.Signature,
		CreatedAt:  data.CreatedAt,
		BlockHash:  s.options.BlockHashHex,
	}
	_, err := s.ctx.DBMongo().Create(models.Transaction{}.TableName(), transaction)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)

	}
	return transaction, nil
}

func (s BCTransactionService) UpdateNonce(didAddress string, signature string) core.IError {
	nonce := models.NewNonce(signature, didAddress)
	_, err := s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address": didAddress,
	}, bson.M{
		"$set": bson.M{
			"nonce": nonce,
		},
	})
	if err != nil {
		return s.ctx.NewError(err, emsgs.DBError)
	}
	return nil
}
func (s BCTransactionService) UpdateDIDDocumentVersion(signature string, didAddress string, createdAt *time.Time) core.IError {
	didDocumentVersion := models.NewDIDDocumentVersion(signature, didAddress, createdAt)
	did, ierr := s.didService.Find(didAddress)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	didHistory := models.DIDHistory{
		Address:    did.Address,
		Controller: did.Controller,
		Keys:       did.Keys,
		Recoverer:  did.Recoverer,
		CreatedAt:  did.CreatedAt,
		UpdatedAt:  did.UpdatedAt,
		Version:    didDocumentVersion,
	}
	_, err := s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address": didAddress,
	}, bson.M{
		"$push": bson.M{
			"history": didHistory,
		},
	})
	if err != nil {
		return s.ctx.NewError(err, emsgs.DBError)
	}
	return nil
}
