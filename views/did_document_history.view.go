package views

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type DIDDocumentHistory struct {
	DIDDocument []DIDDocumentVersion `json:"did_document"`
}

func NewDIDDocumentHistory(did *models.DID) *DIDDocumentHistory {
	didDocuments := make([]DIDDocumentVersion, 0)

	for _, didHistory := range did.History {
		publicKeys := make([]DIDDocumentVerificationMethod, 0)
		for _, key := range didHistory.Keys {
			if key.RevokedAt == nil {
				publicKey := &DIDDocumentVerificationMethod{
					ID:           key.ID,
					Controller:   key.Controller,
					PublicKeyPem: key.PublicKey,
					Type:         key.KeyType,
				}
				utils.Copy(publicKey, key)
				publicKeys = append(publicKeys, *publicKey)
			}
		}
		didDocument := &DIDDocumentVersion{
			Context:            consts.ContextDIDDocument,
			ID:                 did.Address,
			VerificationMethod: publicKeys,
			Controller:         did.Controller,
			Recoverer:          didHistory.Recoverer,
			Version:            didHistory.Version,
		}
		didDocuments = append(didDocuments, *didDocument)

	}
	return &DIDDocumentHistory{
		DIDDocument: didDocuments,
	}
}
