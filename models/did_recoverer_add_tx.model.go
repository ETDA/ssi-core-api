package models

import "time"

type DIDRecovererAddMessage struct {
	Operation  string `json:"operation"`
	DIDAddress string `json:"did_address"`
	Recoverer  string `json:"recoverer"`
	Nonce      string `json:"nonce"`
}

type DIDRecovererAddTX struct {
	TX        []byte                  `json:"tx" bson:"tx"`
	Message   string                  `json:"message" bson:"message"`
	Payload   *DIDRecovererAddMessage `json:"payload" bson:"payload"`
	Signature string                  `json:"signature" bson:"signature"`
	Version   string                  `json:"version" bson:"version"`
	CreatedAt *time.Time              `json:"created_at" bson:"created_at"`
}
