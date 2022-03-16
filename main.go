package main

import (
	"fmt"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/abci"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/consts"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/did"
	"ssi-gitlab.teda.th/finema/etda/ssi-core-api/vcstatus"
)

func main() {
	switch core.NewEnv().Config().Service {
	case string(consts.ServiceDID):
		did.Run()
	case string(consts.ServiceABCI):
		abci.Run()
	case string(consts.ServiceVCStatus):
		vcstatus.Run()

	default:
		fmt.Printf("Service not found")
	}
}
