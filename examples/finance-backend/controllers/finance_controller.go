// finance_controller.go
package controllers

import (
	"finance-backend/models"
	"finance-backend/utils"
	"fmt"
	"strconv"

	"github.com/zlorgoncho1/sprint/core"
)

// FinanceController returns the routes for managing finances (addIncome, addExpense, getSummary)
func FinanceController() *core.Controller {
	var financeController = &core.Controller{Name: "FinanceController", Path: "finance"}
	financeController.AddRoute(core.POST, "addIncome", addIncome)
	financeController.AddRoute(core.POST, "addExpense", addExpense)
	financeController.AddRoute(core.GET, "getSummary", getSummary)
	return financeController
}

func addIncome(request core.Request) core.Response {
	params, err := utils.InterfaceToJSONObj(request.Body)

	if err != nil {
		return core.Response{Content: "Error parsing request body\n"}
	}

	_amount, ok := params["amount"].(string)
	if !ok {
		return core.Response{Content: "Invalid amount\n"}
	}

	amount, err := strconv.ParseFloat(_amount, 64)
	if err != nil {
		return core.Response{Content: "Error parsing amount\n"}
	}

	category, ok := params["category"].(string)
	if !ok {
		return core.Response{Content: "Invalid category\n"}
	}

	err = models.AddIncome(amount, category)
	if err != nil {
		return core.Response{Content: fmt.Sprintf("Error adding income: %v\n", err)}
	}

	return core.Response{Content: "Income added successfully\n"}
}

func addExpense(request core.Request) core.Response {
	params, err := utils.InterfaceToJSONObj(request.Body)
	if err != nil {
		return core.Response{Content: "Error parsing request body\n"}
	}

	_amount, ok := params["amount"].(string)
	if !ok {
		return core.Response{Content: "Invalid amount\n"}
	}

	amount, err := strconv.ParseFloat(_amount, 64)
	if err != nil {
		return core.Response{Content: "Error parsing amount\n"}
	}

	category, ok := params["category"].(string)
	if !ok {
		return core.Response{Content: "Invalid category\n"}
	}

	err = models.AddExpense(amount, category)
	if err != nil {
		return core.Response{Content: fmt.Sprintf("Error adding expense: %v\n", err)}
	}

	return core.Response{Content: "Expense added successfully\n"}
}

func getSummary(request core.Request) core.Response {
	summary := models.GetSummary()
	return core.Response{Content: summary, ContentType: core.JSON}
}
