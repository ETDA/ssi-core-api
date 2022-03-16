package models

import "time"

type DIDKeyRevokeMessage struct {
	CurrentKey string `json:"current_key" bson:"current_key"`
	Operation  string `json:"operation" bson:"operation"`
	DIDAddress string `json:"did_address" bson:"did_address"`
	KeyID      string `json:"key_id" bson:"key_id"`
	Nonce      string `json:"nonce" bson:"nonce"`
}

type DIDKeyRevokeTX struct {
	TX        []byte               `json:"tx" bson:"tx"`
	Message   string               `json:"message" bson:"message"`
	Payload   *DIDKeyRevokeMessage `json:"payload" bson:"payload"`
	Signature string               `json:"signature" bson:"signature"`
	Version   string               `json:"version" bson:"version"`
	CreatedAt *time.Time           `json:"created_at" bson:"created_at"`
}
