package views

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type DIDDocument struct {
	Context            string                          `json:"@context"`
	ID                 string                          `json:"id"`
	Controller         string                          `json:"controller"`
	VerificationMethod []DIDDocumentVerificationMethod `json:"verificationMethod"`
	Recoverer          *string                         `json:"recoverer"`
}

type DIDDocumentVersion struct {
	Context            string                          `json:"@context"`
	ID                 string                          `json:"id"`
	Controller         string                          `json:"controller"`
	VerificationMethod []DIDDocumentVerificationMethod `json:"verificationMethod"`
	Recoverer          *string                         `json:"recoverer"`
	Version            string                          `json:"version"`
}

type DIDDocumentVerificationMethod struct {
	ID           string `json:"id"`
	Controller   string `json:"controller"`
	PublicKeyPem string `json:"publicKeyPem"`
	Type         string `json:"type"`
}

func NewDIDDocument(did *models.DID) *DIDDocument {
	verificationMethods := make([]DIDDocumentVerificationMethod, 0)
	for _, key := range did.Keys {
		if key.RevokedAt == nil {
			publicKey := &DIDDocumentVerificationMethod{
				ID:           key.ID,
				Controller:   key.Controller,
				PublicKeyPem: key.PublicKey,
				Type:         key.KeyType,
			}
			utils.Copy(publicKey, key)
			verificationMethods = append(verificationMethods, *publicKey)
		}

	}
	return &DIDDocument{
		Context:            consts.ContextDIDDocument,
		ID:                 did.Address,
		Controller:         did.Controller,
		VerificationMethod: verificationMethods,
		Recoverer:          did.Recoverer,
	}
}

func NewDIDDocumentVersion(did *models.DIDHistory) *DIDDocumentVersion {
	verificationmethods := make([]DIDDocumentVerificationMethod, 0)
	for _, key := range did.Keys {
		if key.RevokedAt == nil {
			publicKey := &DIDDocumentVerificationMethod{
				ID:           key.ID,
				Controller:   key.Controller,
				PublicKeyPem: key.PublicKey,
				Type:         key.KeyType,
			}
			utils.Copy(publicKey, key)
			verificationmethods = append(verificationmethods, *publicKey)
		}
	}
	return &DIDDocumentVersion{
		Context:            consts.ContextDIDDocument,
		ID:                 did.Address,
		Controller:         did.Controller,
		VerificationMethod: verificationmethods,
		Recoverer:          did.Recoverer,
		Version:            did.Version,
	}
}
