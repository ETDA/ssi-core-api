package models

import (
	"time"

	"ssi-gitlab.teda.th/ssi/core/utils"
)

type VC struct {
	CID         string     `json:"cid" bson:"cid"`
	DIDAddress  string     `json:"did_address" bson:"did_address"`
	VCHash      string     `json:"vc_hash" bson:"vc_hash"`
	Status      *string    `json:"status" bson:"status"`
	Tags        []string   `json:"tags" bson:"tags"`
	ActivatedAt *time.Time `json:"activated_at" bson:"activated_at"`
	CreatedAt   *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" bson:"updated_at"`
	RevokedAt   *time.Time `json:"revoked_at" bson:"revoked_at"`
	ExpiredAt   *time.Time `json:"expired_at" bson:"expired_at"`
}

func (VC) TableName() string {
	return "vcs"
}

func NewCID(didAddress string, createdAt *time.Time) string {
	return utils.NewSha256(didAddress + createdAt.String())
}
