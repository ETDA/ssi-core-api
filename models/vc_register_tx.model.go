package models

import "time"

type VCRegisterMessage struct {
	Operation  string `json:"operation" bson:"operation"`
	DIDAddress string `json:"did_address" bson:"did_address"`
	Nonce      string `json:"nonce" bson:"nonce"`
}

type VCRegisterTX struct {
	TX        []byte             `json:"tx" bson:"tx"`
	Message   string             `json:"message" bson:"message"`
	Payload   *VCRegisterMessage `json:"payload" bson:"payload"`
	Signature string             `json:"signature" bson:"signature"`
	Version   string             `json:"version" bson:"version"`
	CreatedAt *time.Time         `json:"created_at" bson:"created_at"`
}
