package services

import (
	"encoding/json"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/emsgs"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
)

type ValidateVCSchemaSchemaPayload struct {
	ID       string
	Document *json.RawMessage
	VCTypes  []string
}

type ISchemaService interface {
	Validate(payload *ValidateVCSchemaSchemaPayload) (bool, core.Map, core.IError)
}

type schemaService struct {
	ctx core.IContext
}

func NewVCSchemaService(ctx core.IContext) ISchemaService {
	return &schemaService{
		ctx: ctx,
	}
}

type validateResult struct {
	Valid  bool     `json:"valid"`
	Fields core.Map `json:"fields"`
}

func (s schemaService) Validate(payload *ValidateVCSchemaSchemaPayload) (bool, core.Map, core.IError) {
	cc := s.ctx.(core.IHTTPContext)
	header := cc.Request().Header
	header.Set("content-type", "application/json")

	var schemaMapper map[string]interface{}
	ierr := core.RequesterToStruct(&schemaMapper, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(payload.ID, &core.RequesterOptions{
			Headers: header,
		})
	})

	if _, ok := schemaMapper["$schema"]; !ok {
		return false, nil, s.ctx.NewError(emsgs.JSONSchemaNotFound(payload.ID), emsgs.JSONSchemaNotFound(payload.ID))
	}

	schema := json.RawMessage(utils.JSONToString(schemaMapper))

	res := &validateResult{}
	ierr = core.RequesterToStruct(res, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			"/schema/validate-document",
			core.Map{
				"schema":   &schema,
				"document": payload.Document,
			},
			&core.RequesterOptions{
				BaseURL: s.ctx.ENV().String(consts.ENVJSONSchemaValidatorAPIEndpoint),
				Headers: header,
			},
		)
	})

	if ierr != nil {
		return false, nil, s.ctx.NewError(ierr, ierr)
	}

	return res.Valid, res.Fields, nil
}
