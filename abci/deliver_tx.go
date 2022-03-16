package abci

import (
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/middlewares"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/services"
)

func (app *Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	if err := middlewares.ValidateSignatureMessage(app.ctx, req.Tx); err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeUnauthorized, Info: err.Error()}
	}

	err := app.ValidateTx(req.Tx)
	if err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeUnauthorized, Info: err.Error()}
	}

	bcSvc := services.NewBCTransactionService(app.ctx, &services.BCTranasctionServiceOptions{
		BlockHashHex: app.state.GetBlockHashHex(),
	})

	var res *types.ResponseDeliverTx

	switch app.ctx.GetOperation(req.Tx) {
	case consts.OperationDIDRegister:
		payload := &models.DIDRegisterTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}

		did, ierr := bcSvc.RegisterDID(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break
	case consts.OperationDIDRecovererAdd:
		payload := &models.DIDRecovererAddTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}

		did, ierr := bcSvc.DIDAddRecoverer(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break
	case consts.OperationDIDKeyAdd:
		payload := &models.DIDKeyAddTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}

		did, ierr := bcSvc.AddKey(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break

	case consts.OperationDIDKeyReset:
		payload := &models.DIDKeyResetTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}

		did, ierr := bcSvc.ResetKey(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break

	case consts.OperationDIDKeyRevoke:
		payload := &models.DIDKeyRevokeTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}

		did, ierr := bcSvc.RevokeKey(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break

	case consts.OperationVCRegister:
		payload := &models.VCRegisterTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}
		did, ierr := bcSvc.VCRegister(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break
	case consts.OperationVCAddStatus:
		payload := &models.VCUpdateStatusTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}
		did, ierr := bcSvc.VCAddStatus(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break
	case consts.OperationVCUpdateStatus:
		payload := &models.VCUpdateStatusTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}
		did, ierr := bcSvc.VCUpdateStatus(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}

		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK, Info: utils.JSONToString(did)}
		break
	case consts.OperationVCTagStatus:
		payload := &models.VCTagStatusTX{}
		err := app.BindTX(payload, req.Tx)
		if err != nil {
			app.ctx.NewError(err, emsgs.InternalServerError)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: err.Error()}
		}
		_, ierr := bcSvc.VCTagStatus(payload)
		if ierr != nil {
			app.ctx.NewError(ierr, ierr)
			return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: ierr.Error()}
		}
		res = &types.ResponseDeliverTx{Code: code.CodeTypeOK}
	}

	if res != nil {
		app.state.IncreaseSizeOne()
		return *res
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeUnknownError, Info: "SOMETHING WENT WRONG"}
}

func (app *Application) BindTX(dest interface{}, tx []byte) error {
	payload := make(map[string]interface{})
	err := utils.JSONParse(tx, &payload)
	if err != nil {
		return err
	}

	payload["tx"] = tx

	// extract message base64 into payload.Payload
	msgJSON, err := utils.Base64Decode(payload["message"].(string))
	if err != nil {
		return err
	}

	tmpPayload := make(map[string]interface{})
	err = utils.JSONParse(utils.StringToBytes(msgJSON), &tmpPayload)
	if err != nil {
		return err
	}

	payload["payload"] = tmpPayload

	err = utils.MapToStruct(payload, &dest)
	if err != nil {
		return err
	}

	return nil
}
