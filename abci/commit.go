package abci

import (
	"encoding/binary"
	"github.com/tendermint/tendermint/abci/types"
	"log"
)

func (app *Application) Commit() (resp types.ResponseCommit) {
	log.Println("commit...")
	app.state.IncreaseHeightOne()
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.GetSize())
	app.state.setAppHash(appHash)
	return types.ResponseCommit{Data: app.state.getAppHash()}
}
