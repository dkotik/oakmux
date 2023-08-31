package oakmux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrailingSlashRedirects(t *testing.T) {
	mux, err := New(
		WithPrefix("api/v1/"),
		WithRouteHandler("test", "test", newTestHandler(t)),
		WithRouteHandler("test2", "test2/", newTestHandler(t)),
	)
	if err != nil {
		t.Fatal(err)
	}

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test", nil),
		http.StatusOK, "/api/v1/test")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test/", nil),
		http.StatusTemporaryRedirect, "")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test2/", nil),
		http.StatusOK, "/api/v1/test2/")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test2", nil),
		http.StatusTemporaryRedirect, "")(t)
}

func TestTrailingToSlashRedirects(t *testing.T) {
	mux, err := New(
		WithoutTrailingSlashRedirectsFromSlash(),
		WithPrefix("api/v1/"),
		WithRouteHandler("test", "test", newTestHandler(t)),
		WithRouteHandler("test2", "test2/", newTestHandler(t)),
	)
	if err != nil {
		t.Fatal(err)
	}

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test", nil),
		http.StatusOK, "/api/v1/test")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test/", nil),
		http.StatusTemporaryRedirect, "")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test2/", nil),
		http.StatusOK, "/api/v1/test2/")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test2", nil),
		http.StatusNotFound, "")(t)
}

func TestTrailingFromSlashRedirects(t *testing.T) {
	mux, err := New(
		WithoutTrailingSlashRedirectsToSlash(),
		WithPrefix("api/v1/"),
		WithRouteHandler("test", "test", newTestHandler(t)),
		WithRouteHandler("test2", "test2/", newTestHandler(t)),
	)
	if err != nil {
		t.Fatal(err)
	}

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test", nil),
		http.StatusOK, "/api/v1/test")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test/", nil),
		http.StatusNotFound, "")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test2/", nil),
		http.StatusOK, "/api/v1/test2/")(t)

	expectFromRequest(mux,
		httptest.NewRequest(http.MethodPost, "/api/v1/test2", nil),
		http.StatusTemporaryRedirect, "")(t)
}
