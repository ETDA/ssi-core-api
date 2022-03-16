package abci

import (
	"github.com/tendermint/tendermint/abci/types"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"log"
)

func (app *Application) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	log.Println("initChain....")
	app.ctx.DBMongo().Drop((&models.DID{}).TableName())
	app.ctx.DBMongo().Drop((&models.Transaction{}).TableName())
	app.ctx.DBMongo().Drop((&models.VC{}).TableName())
	app.ctx.DBMongo().Drop((&models.Block{}).TableName())
	return types.ResponseInitChain{}
}
