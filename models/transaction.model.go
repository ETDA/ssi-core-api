package models

import "time"

type Transaction struct {
	ID               string     `json:"id" bson:"_id"`
	DIDAddress       string     `json:"did_address" bson:"did_address"`
	KeyID            string     `json:"key_id" bson:"key_id"`
	Message          string     `json:"message" bson:"message"`
	Operation        string     `json:"operation" bson:"operation"`
	Signature        string     `json:"signature" bson:"signature"`
	CreatedAt        *time.Time `json:"created_at" bson:"created_at"`
	BlockHash        string     `json:"block_hash" bson:"block_hash"`
	ConfirmedBlockAt *time.Time `json:"confirmed_block_at" bson:"confirmed_block_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
