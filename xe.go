package main

import (
	"github.com/jeffreyyong/xe/calculator"
	"github.com/jeffreyyong/xe/client"
	"github.com/jeffreyyong/xe/server"
)

func main() {
	addr := "localhost:3000"
	c := client.NewHTTPClient()
	fx := client.NewForex(c)
	ce := calculator.NewEngine()
	h := server.NewHandler(fx, ce)
	httpHandler := server.SetupAPIHandler(h)
	xeService := server.NewXEService(httpHandler, addr)
	xeService.Run()
}
