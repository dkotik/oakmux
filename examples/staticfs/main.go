/*
Package main demonstrates the basic use of [staticfs.FS]
*/
package main

import (
	"embed"
	"fmt"
	"net"
	"net/http"

	"github.com/dkotik/oakmux/staticfs"
)

//go:embed main.go
var fs embed.FS

func main() {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	handler, err := staticfs.New(
		staticfs.WithFileSystem(fs),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("File system contents:", handler.String())
	fmt.Printf("Listening at http://%s/main.go\n", l.Addr())
	http.Serve(l, handler)
}
