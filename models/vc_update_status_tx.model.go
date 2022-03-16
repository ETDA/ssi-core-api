package models

import "time"

type VCUpdateStatusMessage struct {
	Operation  string `json:"operation" bson:"operation"`
	DIDAddress string `json:"did_address" bson:"did_address"`
	CID        string `json:"cid" bson:"cid"`
	VCHash     string `json:"vc_hash" bson:"vc_hash"`
	Status     string `json:"status" bson:"status"`
	Nonce      string `json:"nonce" bson:"nonce"`
}

type VCUpdateStatusTX struct {
	TX        []byte                 `json:"tx" bson:"tx"`
	Message   string                 `json:"message" bson:"message"`
	Payload   *VCUpdateStatusMessage `json:"payload" bson:"payload"`
	Signature string                 `json:"signature" bson:"signature"`
	Version   string                 `json:"version" bson:"version"`
	CreatedAt *time.Time             `json:"created_at" bson:"created_at"`
}
