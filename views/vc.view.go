package views

import (
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type VC struct {
	CID        string `json:"cid"`
	DIDAddress string `json:"did_address"`
}

func NewVC(vc *models.VC) *VC {
	return &VC{
		CID:        vc.CID,
		DIDAddress: vc.DIDAddress,
	}
}
