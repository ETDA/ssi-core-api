package abci

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/requests"
)

func (app *Application) ValidateTx(tx []byte) error {
	msgJSON := app.ctx.GetMessageJSON(tx)

	switch app.ctx.GetOperation(tx) {
	case consts.OperationDIDRegister:
		input := &requests.RequestDIDRegister{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}

		break

	case consts.OperationDIDRecovererAdd:
		input := &requests.RequestDIDRecovererAdd{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}

		break

	case consts.OperationDIDKeyAdd:
		input := &requests.RequestDIDKeyAdd{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}

		break

	case consts.OperationDIDKeyReset:
		input := &requests.RequestDIDKeyReset{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}

		break

	case consts.OperationDIDKeyRevoke:
		input := &requests.RequestDIDRevokeKey{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}

		break

	case consts.OperationVCRegister:
		input := &requests.RequestVCRegister{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}
		break

	case consts.OperationVCAddStatus:
		input := &requests.VCStatusAdd{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}
		break
	case consts.OperationVCUpdateStatus:
		input := &requests.VCStatusUpdate{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}
		break
	case consts.OperationVCTagStatus:
		input := &requests.VCStatusTag{}
		utils.JSONParse(utils.StringToBytes(msgJSON), input)
		if err := input.Valid(app.ctx); err != nil {
			return err
		}
		break
	}
	return nil
}
