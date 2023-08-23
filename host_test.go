package oakmux

import (
	"io"
	"net/http"
	"testing"
)

var testHostHandler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	io.WriteString(w, "Hello world!")
	return nil
})

func TestHostMuxCreation(t *testing.T) {
	mux, err := NewHostMux(
		WithHostHandler("one", testHostHandler),
		WithHostHandler("two", testHostHandler),
		WithHostHandler("three", testHostHandler),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := mux.(*listHostMux)
	if !ok {
		t.Fatal("created mux is not a listHostMux")
	}

	mux, err = NewHostMux(
		WithHostHandler("one", testHostHandler),
		WithHostHandler("two", testHostHandler),
		WithHostHandler("three", testHostHandler),
		WithHostHandler("four", testHostHandler),
		WithHostHandler("five", testHostHandler),
		WithHostHandler("six", testHostHandler),
		WithHostHandler("seven", testHostHandler),
		WithHostHandler("eight", testHostHandler),
		WithHostHandler("nine", testHostHandler),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, ok = mux.(mapHostMux)
	if !ok {
		t.Fatal("created mux is not a mapHostMux")
	}
}
