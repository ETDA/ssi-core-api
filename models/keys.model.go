package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type Key struct {
	ID         string     `json:"id" bson:"_id"`
	Controller string     `json:"controller" bson:"controller"`
	PublicKey  string     `json:"public_key" bson:"public_key"`
	KeyType    string     `json:"key_type" bson:"key_type"`
	CreatedAt  *time.Time `json:"created_at" bson:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at" bson:"last_used_at"`
	RevokedAt  *time.Time `json:"revoked_at" bson:"revoked_at"`
}

func (Key) TableName() string {
	return "keys"
}

func NewKeyID(didAddress string, createdAt *time.Time) string {
	return utils.NewSha256(didAddress + createdAt.String())
}
