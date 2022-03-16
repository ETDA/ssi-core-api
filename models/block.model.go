package models

import "time"

type Block struct {
	ID          string     `json:"id" bson:"_id"`
	Height      int64      `json:"height" bson:"height"`
	ConfirmedAt *time.Time `json:"confirmed_at" bson:"confirmed_at"`
}

func (b Block) TableName() string {
	return "blocks"
}
