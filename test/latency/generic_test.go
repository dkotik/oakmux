package latency

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dkotik/oakmux"
)

func requestFactory() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
}

func stringDecoder(r *http.Request) (string, error) {
	return "test", nil
}

func TestGenericAdaptorLatency(t *testing.T) {
	explicit := testing.Benchmark(BenchmarkExplicitHandler).NsPerOp()
	generic := testing.Benchmark(BenchmarkGenericAdaptor).NsPerOp()
	speed := float32(explicit-generic) * 100 / float32(explicit)
	if speed < 0 {
		t.Fatalf("generic adapter turned out to be %.2f%% slower than the explicit one: %dns/op vs %dns/op", -speed, explicit, generic)
	}
	t.Logf("generic adapter is %.2f%% faster than the explicit one: %dns/op vs %dns/op", speed, explicit, generic)
}

func BenchmarkExplicitHandler(b *testing.B) {
	handler, err := oakmux.New(
		oakmux.WithPrefix("api/v1/"),
		oakmux.WithRouteHandler(
			"test", "test",
			oakmux.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) error {
					input, err := stringDecoder(r)
					if err != nil {
						return err
					}
					err = json.NewEncoder(w).Encode(input + "1")
					return err
				},
			),
		),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		if err = handler.ServeHyperText(w, requestFactory()); err != nil {
			b.Fatal("test request failed:", err)
		}
	}
}

func BenchmarkGenericAdaptor(b *testing.B) {
	handler, err := oakmux.New(
		oakmux.WithPrefix("api/v1/"),
		oakmux.WithRouteStringFunc(
			"test", "test",
			func(ctx context.Context, input string) (output string, err error) {
				return input + "1", nil
			},
			stringDecoder,
		),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		if err = handler.ServeHyperText(w, requestFactory()); err != nil {
			b.Fatal("test request failed:", err)
		}
	}
}
