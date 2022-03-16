package vcstatus

import (
	"fmt"
	core "ssi-gitlab.teda.th/ssi/core"
	"os"
)

func Run() {
	env := core.NewEnv()

	mongoDB, err := core.NewDatabaseMongo(env.Config()).Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "MongoDB: %v", err)
		os.Exit(1)
	}
	defer mongoDB.Close()

	e := core.NewHTTPServer(&core.HTTPContextOptions{
		ContextOptions: &core.ContextOptions{
			MongoDB: mongoDB,
			ENV:     env,
		},
	})

	NewVCStatusHTTPHandler(e)
	e.Logger.Fatal(e.Start(env.Config().Host))
}
