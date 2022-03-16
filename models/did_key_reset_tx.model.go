package models

import (
	"time"
)

type DIDKeyResetMessageNewKey struct {
	PublicKey  string `json:"public_key" bson:"public_key"`
	Signature  string `json:"signature" bson:"signature"`
	Controller string `json:"controller" bson:"controller"`

	// Just add this field for future
	KeyType string `json:"key_type" bson:"key_type"`
}

type DIDKeyResetMessage struct {
	Operation  string                    `json:"operation"`
	DIDAddress string                    `json:"did_address"`
	RequestDID string                    `json:"request_did"`
	NewKey     *DIDKeyResetMessageNewKey `json:"new_key"`
	Nonce      string                    `json:"nonce"`
}

type DIDKeyResetTX struct {
	TX        []byte              `json:"tx" bson:"tx"`
	Message   string              `json:"message" bson:"message"`
	Payload   *DIDKeyResetMessage `json:"payload" bson:"payload"`
	Signature string              `json:"signature" bson:"signature"`
	Version   string              `json:"version" bson:"version"`
	CreatedAt *time.Time          `json:"created_at" bson:"created_at"`
}
