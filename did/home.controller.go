package did

import (
	"ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

type HomeController struct{}

func (n *HomeController) Get(c core.IHTTPContext) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello, I'm DID API",
	})
}
