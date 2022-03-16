package models

import "time"

type DIDRegisterMessage struct {
	PublicKey string `json:"public_key" bson:"public_key"`
	Operation string `json:"operation" bson:"operation"`
	KeyType   string `json:"key_type" bson:"key_type"`
}

type DIDRegisterTX struct {
	TX        []byte              `json:"tx" bson:"tx"`
	Message   string              `json:"message" bson:"message"`
	Payload   *DIDRegisterMessage `json:"payload" bson:"payload"`
	Signature string              `json:"signature" bson:"signature"`
	Version   string              `json:"version" bson:"version"`
	CreatedAt *time.Time          `json:"created_at" bson:"created_at"`
}
