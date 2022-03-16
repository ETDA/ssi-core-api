package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type DID struct {
	Address    string       `json:"address" bson:"address"`
	Controller string       `json:"controller" bson:"controller"`
	Recoverer  *string      `json:"recoverer" bson:"recoverer"`
	Nonce      string       `json:"nonce" bson:"nonce"`
	Keys       []Key        `json:"keys" bson:"keys"`
	History    []DIDHistory `json:"history" bson:"history"`
	CreatedAt  *time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time   `json:"updated_at" bson:"updated_at"`
}

type DIDHistory struct {
	Address    string     `json:"address" bson:"address"`
	Controller string     `json:"controller" bson:"controller"`
	Recoverer  *string    `json:"recoverer" bson:"recoverer"`
	Keys       []Key      `json:"keys" bson:"keys"`
	CreatedAt  *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at" bson:"updated_at"`
	Version    string     `json:"version" bson:"version"`
}

func (DID) TableName() string {
	return "dids"
}

func NewNonce(signature string, did string) string {
	return utils.GetMD5Hash(signature + did)
}

func NewDIDAddress(currentPublicKey string, didMethod string) string {
	return utils.GenerateDID(utils.NewSha256(currentPublicKey), didMethod)
}

func NewDIDDocumentVersion(signature string, did string, createdAt *time.Time) string {
	return utils.GetMD5Hash(signature + did + createdAt.String())
}
