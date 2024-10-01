package module

import (
	"finance-backend/controllers"

	"github.com/zlorgoncho1/sprint/core"
)

func FinanceModule() *core.Module {
	return &core.Module{
		Name: "FinanceModule",
		Controllers: []*core.Controller{
			controllers.UserController(),
			controllers.FinanceController(),
		},
	}
}
