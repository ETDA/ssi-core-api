package services

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"ssi-gitlab.teda.th/ssi/core"
	"go.mongodb.org/mongo-driver/bson"
)

func (s BCTransactionService) RegisterDID(data *models.DIDRegisterTX) (*models.DID, core.IError) {
	did := models.NewDIDAddress(data.Payload.PublicKey, s.ctx.ENV().Config().DIDMethodDefault)
	keyID := models.NewKeyID(did, data.CreatedAt)
	nonce := models.NewNonce(data.Signature, data.Payload.PublicKey)
	didDocumentVersion := models.NewDIDDocumentVersion(data.Signature, did, data.CreatedAt)
	newDID := &models.DID{
		Address:    did,
		Controller: did,
		Nonce:      nonce,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.CreatedAt,
		Keys: []models.Key{{
			ID:         keyID,
			Controller: did,
			CreatedAt:  data.CreatedAt,
			LastUsedAt: data.CreatedAt,
			RevokedAt:  nil,
			PublicKey:  data.Payload.PublicKey,
			KeyType:    data.Payload.KeyType,
		}},
		History: []models.DIDHistory{{
			Address:    did,
			Controller: did,
			Keys: []models.Key{{
				ID:         keyID,
				Controller: did,
				CreatedAt:  data.CreatedAt,
				LastUsedAt: data.CreatedAt,
				RevokedAt:  nil,
				PublicKey:  data.Payload.PublicKey,
				KeyType:    data.Payload.KeyType,
			}},
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.CreatedAt,
			Version:   didDocumentVersion,
		}},
	}
	_, err := s.ctx.DBMongo().Create(newDID.TableName(), newDID)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	_, ierr := s.createTransaction(&createTransactionPayload{
		TX:         data.TX,
		DIDAddress: did,
		KeyID:      keyID,
		Operation:  data.Payload.Operation,
		Signature:  data.Signature,
		CreatedAt:  data.CreatedAt,
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return s.didService.Find(did)
}

func (s BCTransactionService) DIDAddRecoverer(data *models.DIDRecovererAddTX) (*models.DID, core.IError) {
	publicKey, ierr := s.keyService.FindVerifiedPublicKey(data.Payload.DIDAddress, data.Message, data.Signature)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	_, err := s.ctx.DBMongo().UpdateOne(models.DID{}.TableName(), bson.M{
		"address":  data.Payload.DIDAddress,
	}, bson.M{
		"$set": bson.M{
			"recoverer": data.Payload.Recoverer,
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
