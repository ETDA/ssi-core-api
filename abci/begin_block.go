package abci

import (
	"github.com/tendermint/tendermint/abci/types"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
)

func (app *Application) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	bcSvc := services.NewBCTransactionService(app.ctx, &services.BCTranasctionServiceOptions{})

	app.state.SetBlockHash(req.GetHash())
	_, ierr := bcSvc.CreateBlock(app.state.GetBlockHashHex(), req.Header.Height)
	if ierr != nil {
		app.ctx.NewError(ierr, ierr)
	}
	req.Header.LastBlockId.GetHash()
	previousBlockHash := utils.BytesToHexString(req.Header.LastBlockId.GetHash())

	time := req.Header.GetTime()

	ierr = bcSvc.ConfirmBlock(previousBlockHash, &time)

	if ierr != nil {
		app.ctx.NewError(ierr, ierr)
	}

	return types.ResponseBeginBlock{}
}
