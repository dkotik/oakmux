package oakmux

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func expectFromRequest(h Handler, r *http.Request, code int, body string) func(*testing.T) {
	return func(t *testing.T) {
		w := httptest.NewRecorder()
		err := h.ServeHyperText(w, r)
		result := w.Result()
		if result.StatusCode != code {
			var httpError Error
			if err != nil && errors.As(err, &httpError) {
				if httpError.HyperTextStatusCode() != code {
					t.Fatalf("status code does not match expected result for path %q: %d vs %d", r.URL.Path, httpError.HyperTextStatusCode(), code)
				}
			} else {
				t.Fatalf("status code does not match expected result for path %q: %d vs %d", r.URL.Path, result.StatusCode, code)
			}
		}
		var b bytes.Buffer
		_, _ = io.Copy(&b, result.Body)
		_ = result.Body.Close()
		if body != b.String() {
			t.Fatalf("body does not match expected result for path %q: %q vs %q", r.URL.Path, body, b.String())
		}
		if err == nil {
			return
		}
		var httpError Error
		if !errors.As(err, &httpError) {
			t.Fatal("unexpected error occured:", err)
		}
	}
}

func newTestHandler(t *testing.T) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		t.Log("test request suceeded:", r.URL.Path)
		io.WriteString(w, r.URL.Path)
		return nil
	})
}

func TestMux(t *testing.T) {
	handler := newTestHandler(t)
	mux, err := New(
		WithRouteHandler("firstRoute", "/test/[pattern]/yep/", handler),
		WithRouteHandler("secondRoute", "/test/[wild]/[pattern1]/last/", handler),
	)
	if err != nil {
		t.Fatal(err)
	}

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/test/something/yep/", nil),
		http.StatusOK, "/test/something/yep/")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/test/1/2/last/", nil),
		http.StatusOK, "/test/1/2/last/")(t)
}
