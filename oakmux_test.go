package oakmux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux(t *testing.T) {
	mux := New(
		WithRoute(
			"firstRoute",
			"/test/{pattern}/yep/",
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				t.Log("test request suceeded")
				return nil
			}),
		),
		WithRoute(
			"secondRoute",
			"/test/{wild}/{pattern1}/yep/",
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				t.Log("test request suceeded")
				return nil
			}),
		),
	)

	r := httptest.NewRequest(
		http.MethodPost,
		"/test/something/yep/",
		nil,
	)
	w := httptest.NewRecorder()

	if err := mux.ServeHyperText(w, r); err != nil {
		t.Fatal("mux route failed:", err)
	}
}
