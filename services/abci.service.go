package services

import (
	"errors"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"net/url"
	"time"
)

type BroadcastTXResponse struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	Result struct {
		CheckTx struct {
			Code      int           `json:"code"`
			Codespace string        `json:"codespace"`
			Data      interface{}   `json:"data"`
			Events    []interface{} `json:"events"`
			GasUsed   string        `json:"gasUsed"`
			GasWanted string        `json:"gasWanted"`
			Info      string        `json:"info"`
			Log       string        `json:"log"`
		} `json:"check_tx"`
		DeliverTx struct {
			Code      int           `json:"code"`
			Codespace string        `json:"codespace"`
			Data      interface{}   `json:"data"`
			Events    []interface{} `json:"events"`
			GasUsed   string        `json:"gasUsed"`
			GasWanted string        `json:"gasWanted"`
			Info      string        `json:"info"`
			Log       string        `json:"log"`
		} `json:"deliver_tx"`
		Hash   string `json:"hash"`
		Height string `json:"height"`
	} `json:"result"`
}

func NewABCIService(ctx core.IContext) IABCIService {
	return &ABCIService{ctx}
}

type IABCIService interface {
	BroadcastTXCommit(payload *TXBroadcastPayload) (*BroadcastTXResponse, core.IError)
}

type ABCIService struct {
	ctx core.IContext
}

type TXBroadcastPayload struct {
	Message   string     `json:"message"`
	Signature string     `json:"signature"`
	Version   string     `json:"version"`
	CreatedAt *time.Time `json:"created_at"`
}

func (s ABCIService) BroadcastTXCommit(payload *TXBroadcastPayload) (*BroadcastTXResponse, core.IError) {
	hex := utils.StructToHexString(payload)
	res, err := s.ctx.Requester().Get("/broadcast_tx_commit",
		&core.RequesterOptions{
			BaseURL: s.ctx.ENV().Config().ABCIEndpoint,
			Params: url.Values{
				"tx": {"0x" + hex},
			},
		})

	var resBody BroadcastTXResponse
	b := utils.JSONToString(res.Data)
	_ = utils.JSONParse([]byte(b), &resBody)
	if err != nil {
		return &resBody, s.ctx.NewError(err, emsgs.BroadcastTXError)
	}

	if resBody.Error.Code != 0 {
		return &resBody, s.ctx.NewError(
			errors.New(resBody.Error.Data),
			emsgs.BroadcastTXErrorWithMSG(resBody.Error.Data),
			resBody)
	}

	if resBody.Result.DeliverTx.Code != 0 {
		return &resBody, s.ctx.NewError(
			errors.New(resBody.Result.DeliverTx.Info),
			emsgs.BroadcastTXErrorWithMSG(resBody.Result.DeliverTx.Info),
			resBody)
	}

	if resBody.Result.CheckTx.Code != 0 {
		return &resBody, s.ctx.NewError(
			errors.New(resBody.Result.CheckTx.Info),
			emsgs.BroadcastTXErrorWithMSG(resBody.Result.CheckTx.Info),
			resBody)
	}

	return &resBody, nil
}
