package models

import "time"

type DIDKeyAddMessageNewKey struct {
	PublicKey  string `json:"public_key"`
	Signature  string `json:"signature"`
	Controller string `json:"controller"`

	// Just add this field for future
	KeyType string `json:"key_type"`
}

type DIDKeyAddMessage struct {
	Operation  string                  `json:"operation"`
	DIDAddress string                  `json:"did_address"`
	NewKey     *DIDKeyAddMessageNewKey `json:"new_key"`
	Nonce      string                  `json:"nonce"`
}

type DIDKeyAddTX struct {
	TX        []byte            `json:"tx"`
	Message   string            `json:"message"`
	Payload   *DIDKeyAddMessage `json:"payload"`
	Signature string            `json:"signature"`
	Version   string            `json:"version"`
	CreatedAt *time.Time        `json:"created_at"`
}
