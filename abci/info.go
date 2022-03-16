package abci

import (
	"fmt"
	"github.com/tendermint/tendermint/abci/types"
	"log"
)

func (app *Application) Info(req types.RequestInfo) types.ResponseInfo {
	log.Println("info...")
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.GetSize()),
		LastBlockHeight:  app.state.GetHeight(),
		LastBlockAppHash: app.state.GetBlockHash(),
	}
}
