package models

import "time"

type VCTagStatusMessage struct {
	Operation  string   `json:"operation" bson:"operation"`
	DIDAddress string   `json:"did_address" bson:"did_address"`
	CIDs       []string `json:"cids" bson:"cids"`
	Tags       []string `json:"tags", bson:"tags"`
	Nonce      string   `json:"nonce" bson:"nonce"`
}

type VCTagStatusTX struct {
	TX        []byte              `json:"tx" bson:"tx"`
	Message   string              `json:"message" bson:"message"`
	Payload   *VCTagStatusMessage `json:"payload" bson:"payload"`
	Signature string              `json:"signature" bson:"signature"`
	Version   string              `json:"version" bson:"version"`
	CreatedAt *time.Time          `json:"created_at" bson:"created_at"`
}
