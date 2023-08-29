/*
Package main demonstrates routing directly to domain function calls.
*/
package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/dkotik/oakmux"
)

//go:embed main.go
var fs embed.FS

func main() {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	domainLogic := &OnlineStore{}
	handler, err := oakmux.New(
		oakmux.WithRouteFunc( // Unary
			"order", "api/v1/order",
			domainLogic.Order,
		),
		oakmux.WithRouteStringFunc( // UnaryString
			"price", "api/v1/price",
			domainLogic.GetPrice,
			func(r *http.Request) (string, error) {
				// string decoder
				log.Println("got query:", r.URL.RawQuery)
				return r.URL.Query().Get("item"), nil
			},
		),
		oakmux.WithRouteNullaryFunc( // Nullary
			"inventory", "api/v1/inventory",
			domainLogic.GetInventory,
		),
		oakmux.WithRouteVoidFunc( // Unary Void
			"record", "api/v1/record",
			domainLogic.Record,
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		`Listening at http://%[1]s/

    Test Order (Unary):
      curl -v -d '{"item":"box","quantity":1}' -H 'Content-Type: application/json' http://%[1]s/api/v1/order
    Test Price (Unary String):
      curl -v -G -d 'item=shirt' http://%[1]s/api/v1/price
    Test Inventory (Nullary):
      curl -v http://%[1]s/api/v1/inventory
    Test Record (Unary Void):
      curl -v -d '{"item":"box","quantity":1}' -H 'Content-Type: application/json' http://%[1]s/api/v1/record

`,
		l.Addr(),
	)

	http.Serve(l, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := handler.ServeHyperText(w, r)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		},
	))
}
