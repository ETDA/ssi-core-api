package abci

import (
	"github.com/tendermint/tendermint/abci/types"
	"ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type state struct {
	size         int64
	appHash      []byte
	blockHash    []byte
	blockHashHex string
	height       int64
}

func (s *state) setAppHash(hash []byte) {
	s.appHash = hash
}

func (s state) getAppHash() []byte {
	return s.appHash
}

func (s state) GetSize() int64 {
	return s.size
}

func (s *state) IncreaseSizeOne() {
	s.size++
}

func (s state) GetHeight() int64 {
	return s.height
}

func (s *state) IncreaseHeightOne() {
	s.height++
}

func (s state) GetBlockHash() []byte {
	return s.blockHash
}

func (s *state) SetBlockHash(hash []byte) {
	s.blockHash = hash
}

func (s state) GetBlockHashHex() string {
	return utils.BytesToHexString(s.blockHash)
}

func NewState() IState {
	return &state{}
}

type IState interface {
	GetSize() int64
	IncreaseSizeOne()
	GetHeight() int64
	IncreaseHeightOne()
	GetBlockHash() []byte
	SetBlockHash(hash []byte)
	GetBlockHashHex() string
	setAppHash(hash []byte)
	getAppHash() []byte
}

type Application struct {
	types.BaseApplication
	ctx   core.IABCIContext
	state IState
}

func NewApplication(ctx core.IABCIContext) *Application {
	return &Application{
		ctx:   ctx,
		state: NewState(),
	}
}

// life cycle
//type Application interface {
//	// Info/Query Connection
//	Info(RequestInfo) ResponseInfo                // Return application info
//	SetOption(RequestSetOption) ResponseSetOption // Set application option
//	Query(RequestQuery) ResponseQuery             // Query for state
//
//	// Mempool Connection
//	CheckTx(RequestCheckTx) ResponseCheckTx // Validate a tx for the mempool
//
//	// Consensus Connection
//	InitChain(RequestInitChain) ResponseInitChain    // Initialize blockchain w validators/other info from TendermintCore
//	BeginBlock(RequestBeginBlock) ResponseBeginBlock // Signals the beginning of a block
//	DeliverTx(RequestDeliverTx) ResponseDeliverTx    // Deliver a tx for full processing
//	EndBlock(RequestEndBlock) ResponseEndBlock       // Signals the end of a block, returns changes to the validator set
//	Commit() ResponseCommit                          // Commit the state and return the application Merkle root hash
//}

func (app *Application) SetOption(req types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}
