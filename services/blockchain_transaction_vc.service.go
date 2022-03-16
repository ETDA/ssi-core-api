package services

import (
	"time"

	core "ssi-gitlab.teda.th/ssi/core"
	"go.mongodb.org/mongo-driver/bson"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

func (s BCTransactionService) VCRegister(data *models.VCRegisterTX) (*models.VC, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	cid := models.NewCID(data.Payload.DIDAddress, data.CreatedAt)
	newVC := &models.VC{
		CID:        cid,
		DIDAddress: data.Payload.DIDAddress,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.CreatedAt,
	}
	_, err := s.ctx.DBMongo().Create(newVC.TableName(), newVC)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, ierr = s.createTransaction(&createTransactionPayload{
		TX:         data.TX,
		DIDAddress: data.Payload.DIDAddress,
		KeyID:      publicKey.ID,
		Operation:  data.Payload.Operation,
		Signature:  data.Signature,
		CreatedAt:  data.CreatedAt,
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return s.vcService.Find(cid)
}

func (s BCTransactionService) VCAddStatus(data *models.VCUpdateStatusTX) (*models.VC, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, err := s.ctx.DBMongo().UpdateOne(models.VC{}.TableName(), bson.M{
		"did_address": data.Payload.DIDAddress,
		"cid":         data.Payload.CID,
	}, bson.M{
		"$set": bson.M{
			"vc_hash":      data.Payload.VCHash,
			"status":       data.Payload.Status,
			"activated_at": data.CreatedAt,
			"updated_at":   data.CreatedAt,
		},
	})
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, ierr = s.createTransaction(&createTransactionPayload{
		TX:         data.TX,
		DIDAddress: data.Payload.DIDAddress,
		KeyID:      publicKey.ID,
		Operation:  data.Payload.Operation,
		Signature:  data.Signature,
		CreatedAt:  data.CreatedAt,
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return s.vcService.Find(data.Payload.CID)
}

func (s BCTransactionService) VCUpdateStatus(data *models.VCUpdateStatusTX) (*models.VC, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	var revokedAt *time.Time = nil
	if data.Payload.Status == consts.VCStatusRevoke {
		revokedAt = data.CreatedAt
	}

	var expiredAt *time.Time = nil
	if data.Payload.Status == consts.VCStatusExpired {
		revokedAt = data.CreatedAt
	}

	_, err := s.ctx.DBMongo().UpdateOne(models.VC{}.TableName(), bson.M{
		"did_address": data.Payload.DIDAddress,
		"cid":         data.Payload.CID,
	}, bson.M{
		"$set": bson.M{
			"status":     data.Payload.Status,
			"updated_at": data.CreatedAt,
			"revoked_at": revokedAt,
			"expired_at": expiredAt,
		},
	})
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, ierr = s.createTransaction(&createTransactionPayload{
		TX:         data.TX,
		DIDAddress: data.Payload.DIDAddress,
		KeyID:      publicKey.ID,
		Operation:  data.Payload.Operation,
		Signature:  data.Signature,
		CreatedAt:  data.CreatedAt,
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return s.vcService.Find(data.Payload.CID)
}

func (s BCTransactionService) VCTagStatus(data *models.VCTagStatusTX) ([]models.VC, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	for _, cid := range data.Payload.CIDs {
		vc, ierr := s.vcService.Find(cid)
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
		tags := data.Payload.Tags
		if len(vc.Tags) != 0 {
			tags = append(tags, vc.Tags...)
		}
		_, err := s.ctx.DBMongo().UpdateOne(models.VC{}.TableName(), bson.M{
			"cid": cid,
		}, bson.M{
			"$set": bson.M{
				"tags":       tags,
				"updated_at": data.CreatedAt,
			},
		})
		if err != nil {
			return nil, s.ctx.NewError(err, emsgs.DBError)
		}

	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, ierr = s.createTransaction(&createTransactionPayload{
		TX:         data.TX,
		DIDAddress: data.Payload.DIDAddress,
		KeyID:      publicKey.ID,
		Operation:  data.Payload.Operation,
		Signature:  data.Signature,
		CreatedAt:  data.CreatedAt,
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	status := make([]models.VC, 0)
	for _, cid := range data.Payload.CIDs {
		vc, _ := s.vcService.Find(cid)
		status = append(status, *vc)
	}
	return nil, nil
}
