package abci

import (
	"github.com/tendermint/tendermint/abci/types"
)

func (app *Application) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{}
}
