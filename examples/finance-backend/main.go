package main

import (
	"finance-backend/models"
	"finance-backend/module"

	"github.com/zlorgoncho1/sprint/server"
)

func main() {
	models.InitDatabase()
	srv := &server.Server{Host: "localhost", Port: "8000"}

	srv.Use(server.CORSMiddleware)

	srv.Start(module.FinanceModule())
}
