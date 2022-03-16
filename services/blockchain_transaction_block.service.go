package services

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"ssi-gitlab.teda.th/ssi/core"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (s BCTransactionService) CreateBlock(hash string, height int64) (*models.Block, core.IError) {
	block := &models.Block{
		ID:          hash,
		Height:      height,
		ConfirmedAt: nil,
	}
	_, err := s.ctx.DBMongo().Create(models.Block{}.TableName(), block)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.DBError)
	}

	return block, nil
}

func (s BCTransactionService) ConfirmBlock(hash string, confirmedAt *time.Time) core.IError {

	_, err := s.ctx.DBMongo().UpdateOne(models.Block{}.TableName(), bson.M{
		"_id": hash,
	}, core.NewMongoHelper().Set(bson.M{
		"confirmed_at": confirmedAt,
	}))

	if err != nil {
		return s.ctx.NewError(err, emsgs.DBError)
	}

	_, err = s.ctx.DBMongo().UpdateOne(models.Transaction{}.TableName(), bson.M{
		"block_hash": hash,
	}, core.NewMongoHelper().Set(bson.M{
		"confirmed_block_at": confirmedAt,
	}))

	if err != nil {
		return s.ctx.NewError(err, emsgs.DBError)
	}
	return nil
}
