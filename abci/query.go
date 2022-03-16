package abci

import (
	"fmt"
	"github.com/tendermint/tendermint/abci/types"
	"log"
	"strconv"
)

func (app *Application) Query(req types.RequestQuery) types.ResponseQuery {
	log.Println("queryTX...")

	switch req.Path {
	case "size":
		return types.ResponseQuery{
			Value: []byte(fmt.Sprintf("%v", app.state.GetSize())),
			Log:   strconv.FormatInt(app.state.GetSize(), 10),
		}
	case "height":
		return types.ResponseQuery{
			Value: []byte(fmt.Sprintf("%v", app.state.GetHeight())),
			Log:   strconv.FormatInt(app.state.GetHeight(), 10),
		}
	default:
		return types.ResponseQuery{Log: fmt.Sprintf("Invalid query path. Expected hash or tx, got %v", req.Path)}
	}
}
