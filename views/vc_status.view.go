package views

import (
	"time"

	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/models"
)

type VCStatus struct {
	CID         string     `json:"cid"`
	DIDAddress  string     `json:"did_address"`
	Status      *string    `json:"status"`
	VCHash      string     `json:"vc_hash"`
	Tags        []string   `json:"tags"`
	ActivatedAt *time.Time `json:"activated_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}

func NewVCStatus(vc *models.VC) *VCStatus {
	return &VCStatus{
		CID:         vc.CID,
		DIDAddress:  vc.DIDAddress,
		VCHash:      vc.VCHash,
		Tags:        vc.Tags,
		Status:      vc.Status,
		ActivatedAt: vc.ActivatedAt,
		RevokedAt:   vc.RevokedAt,
	}
}

func NewVCStatusList(vcs []models.VC) []VCStatus {
	newVCS := make([]VCStatus, 0)
	for _, vc := range vcs {
		newVCS = append(newVCS, *NewVCStatus(&vc))
	}
	return newVCS
}
