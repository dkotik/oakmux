/*
Package main demonstrates routing directly to domain function calls.
*/
package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/dkotik/oakmux"
)

func writeMethodUsed(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "got a request using method %q\n", r.Method)
	return nil
}

func main() {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	handler, err := oakmux.New(
		oakmux.WithPrefix("api/v1/"),
		oakmux.WithRouteHandler("order", "order",
			oakmux.Must(oakmux.NewMethodMux(
				oakmux.WithGetHandler(oakmux.HandlerFunc(writeMethodUsed)),
				oakmux.WithPostHandler(oakmux.HandlerFunc(writeMethodUsed)),
				oakmux.WithPutHandler(oakmux.HandlerFunc(writeMethodUsed)),
				oakmux.WithPatchHandler(oakmux.HandlerFunc(writeMethodUsed)),
				oakmux.WithDeleteHandler(oakmux.HandlerFunc(writeMethodUsed)),
			)),
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		`Listening at http://%[1]s/

    Test Get:
      curl -v http://%[1]s/api/v1/order
    Test Post:
      curl -v -X POST http://%[1]s/api/v1/order
    Test Put:
      curl -v -X PUT http://%[1]s/api/v1/order
    Test Patch:
      curl -v -X PATCH http://%[1]s/api/v1/order
    Test Delete:
      curl -v -X DELETE http://%[1]s/api/v1/order
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
