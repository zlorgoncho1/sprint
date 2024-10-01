package controllers

import (
	"finance-backend/models"
	"finance-backend/utils"
	"net/http"

	"github.com/zlorgoncho1/sprint/core"
)

// UserController returns the routes for user management (register, login, profile)
func UserController() *core.Controller {
	var userController = &core.Controller{Name: "UserController", Path: "user"}
	userController.AddRoute(core.POST, "register", registerUser)
	userController.AddRoute(core.POST, "login", loginUser)
	userController.AddRoute(core.GET, "profile", getUserProfile)
	return userController
}

func registerUser(request core.Request) core.Response {
	params, err := utils.InterfaceToJSONObj(request.Body)
	if err != nil {
		return core.Response{Content: "Error parsing request body"}
	}

	username, ok := params["username"].(string)
	if !ok {
		return core.Response{Content: "Invalid username"}
	}

	password, ok := params["password"].(string)
	if !ok {
		return core.Response{Content: "Invalid password"}
	}

	models.RegisterUser(username, password)

	return core.Response{Content: "User registration failed", StatusCode: http.StatusBadRequest}
}

func loginUser(request core.Request) core.Response {
	params, err := utils.InterfaceToJSONObj(request.Body)
	if err != nil {
		return core.Response{Content: "Error parsing request body"}
	}

	username, ok := params["username"].(string)
	if !ok {
		return core.Response{Content: "Invalid username"}
	}

	password, ok := params["password"].(string)
	if !ok {
		return core.Response{Content: "Invalid password"}
	}

	if models.AuthenticateUser(username, password) {
		return core.Response{Content: "Login successful", StatusCode: http.StatusOK}
	}
	return core.Response{Content: "Invalid username or password", StatusCode: http.StatusUnauthorized}
}

func getUserProfile(request core.Request) core.Response {
	params, err := utils.InterfaceToJSONObj(request.Body)
	if err != nil {
		return core.Response{Content: "Error parsing request body"}
	}

	username, ok := params["username"].(string)
	if !ok {
		return core.Response{Content: "Invalid username"}
	}

	user := models.GetUserProfile(username)
	if user != nil {
		return core.Response{Content: user, ContentType: core.JSON}
	}
	return core.Response{Content: "User not found", StatusCode: http.StatusNotFound}
}
