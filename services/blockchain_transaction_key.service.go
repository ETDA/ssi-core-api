package services

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

func (s BCTransactionService) RevokeKey(data *models.DIDKeyRevokeTX) (*models.DID, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, err := s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address":  data.Payload.DIDAddress,
		"keys._id": data.Payload.KeyID,
	}, bson.M{
		"$set": bson.M{
			"keys.$[ele1].revoked_at": data.CreatedAt,
		},
	}, &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"ele1._id": data.Payload.KeyID},
			},
		},
	})

	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	ierr = s.UpdateDIDDocumentVersion(data.Signature, data.Payload.DIDAddress, data.CreatedAt)
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
	return s.didService.Find(data.Payload.DIDAddress)
}

func (s BCTransactionService) AddKey(data *models.DIDKeyAddTX) (*models.DID, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	keyID := models.NewKeyID(data.Payload.DIDAddress, data.CreatedAt)

	controller := data.Payload.DIDAddress
	if data.Payload.NewKey.Controller != "" {
		controller = data.Payload.NewKey.Controller
	}

	keyType := consts.KeyTypeSecp256r12018
	if data.Payload.NewKey.KeyType != "" {
		keyType = data.Payload.NewKey.KeyType
	}

	newKey := &models.Key{
		ID:         keyID,
		KeyType:    keyType,
		Controller: controller,
		CreatedAt:  data.CreatedAt,
		LastUsedAt: data.CreatedAt,
		RevokedAt:  nil,
		PublicKey:  data.Payload.NewKey.PublicKey,
	}

	_, err := s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address": data.Payload.DIDAddress,
	}, bson.M{
		"$push": bson.M{
			"keys": newKey,
		},
	})

	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	ierr = s.UpdateDIDDocumentVersion(data.Signature, data.Payload.DIDAddress, data.CreatedAt)
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

	return s.didService.Find(data.Payload.DIDAddress)
}

func (s BCTransactionService) ResetKey(data *models.DIDKeyResetTX) (*models.DID, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	keyID := models.NewKeyID(data.Payload.RequestDID, data.CreatedAt)

	controller := data.Payload.RequestDID
	if data.Payload.NewKey.Controller != "" {
		controller = data.Payload.NewKey.Controller
	}

	keyType := consts.KeyTypeSecp256r12018
	if data.Payload.NewKey.KeyType != "" {
		keyType = data.Payload.NewKey.KeyType
	}

	newKey := &models.Key{
		ID:         keyID,
		KeyType:    keyType,
		Controller: controller,
		CreatedAt:  data.CreatedAt,
		LastUsedAt: data.CreatedAt,
		RevokedAt:  nil,
		PublicKey:  data.Payload.NewKey.PublicKey,
	}

	_, err := s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address":   data.Payload.RequestDID,
		"recoverer": data.Payload.DIDAddress,
	}, bson.M{
		"$set": bson.M{
			"keys.$[ele1].revoked_at": data.CreatedAt,
		},
	}, &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.M{
					"ele1.revoked_at": nil,
				},
			},
		},
	})
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	_, err = s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address":   data.Payload.RequestDID,
		"recoverer": data.Payload.DIDAddress,
	}, bson.M{
		"$push": bson.M{
			"keys": newKey,
		},
	})
	_, err = s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address": data.Payload.RequestDID,
	}, bson.M{
		"$set": bson.M{
			"recoverer": nil,
		},
	})
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	ierr = s.UpdateNonce(data.Payload.DIDAddress, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	ierr = s.UpdateDIDDocumentVersion(data.Payload.NewKey.Signature, data.Payload.RequestDID, data.CreatedAt)
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

	return s.didService.Find(data.Payload.RequestDID)
}
