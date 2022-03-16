package abci

import (
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/middlewares"
)

func (app *Application) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	if err := middlewares.ValidateSignatureMessage(app.ctx, req.Tx); err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeUnauthorized, Info: err.Error()}
	}
	err := app.ValidateTx(req.Tx)
	if err != nil {
		return types.ResponseCheckTx{Code: code.CodeTypeUnauthorized, Info: err.Error()}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}
