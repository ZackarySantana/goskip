package skip_reverse_proxy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	skip_reverse_proxy "github.com/zackarysantana/goskip/reverse_proxy"
)

func TestReverseProxy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("HandlesRequest", func(t *testing.T) {
		rp := skip_reverse_proxy.New(&url.URL{Scheme: "http", Host: "localhost:8080", Path: "/v1/%s"})
		require.NotNil(t, rp)

		t.Run("NoServer", func(t *testing.T) {
			request, err := http.NewRequestWithContext(ctx, "GET", "/streams/1234", nil)
			require.NoError(t, err)
			require.NotNil(t, request)

			w := httptest.NewRecorder()
			rp.ServeHTTP(w, request)

			resp := w.Result()
			assert.Equal(t, 502, resp.StatusCode)
		})

		t.Run("WithDefaultSettings", func(t *testing.T) {
			var path string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				w.WriteHeader(200)
			}))
			defer srv.Close()

			rp := skip_reverse_proxy.New(&url.URL{
				Scheme: "http",
				Host:   srv.URL[7:],
				Path:   "/v1/%s",
			})
			require.NotNil(t, rp)

			t.Run("FirstRequest", func(t *testing.T) {
				request, err := http.NewRequestWithContext(ctx, "GET", "/streams/1234", nil)
				require.NoError(t, err)
				require.NotNil(t, request)

				w := httptest.NewRecorder()
				rp.ServeHTTP(w, request)

				resp := w.Result()
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "/v1/1234", path)
			})

			t.Run("SecondRequest", func(t *testing.T) {
				request2, err := http.NewRequestWithContext(ctx, "GET", "/streams/5678", nil)
				require.NoError(t, err)
				require.NotNil(t, request2)

				w := httptest.NewRecorder()
				rp.ServeHTTP(w, request2)

				resp := w.Result()
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "/v1/5678", path)
			})
		})

		t.Run("WithCustomPath", func(t *testing.T) {
			var path string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				w.WriteHeader(200)
			}))
			defer srv.Close()

			rp := skip_reverse_proxy.New(&url.URL{
				Scheme: "http",
				Host:   srv.URL[7:],
				Path:   "/%s/unique_path",
			})
			require.NotNil(t, rp)

			t.Run("FirstRequest", func(t *testing.T) {
				request, err := http.NewRequestWithContext(ctx, "GET", "/streams/1234", nil)
				require.NoError(t, err)
				require.NotNil(t, request)

				w := httptest.NewRecorder()
				rp.ServeHTTP(w, request)

				resp := w.Result()
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "/1234/unique_path", path)
			})

			t.Run("SecondRequest", func(t *testing.T) {
				request2, err := http.NewRequestWithContext(ctx, "GET", "/streams/5678", nil)
				require.NoError(t, err)
				require.NotNil(t, request2)

				w := httptest.NewRecorder()
				rp.ServeHTTP(w, request2)

				resp := w.Result()
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "/5678/unique_path", path)
			})
		})

		t.Run("WithCustomGetUUIDFromPath", func(t *testing.T) {
			var path string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				w.WriteHeader(200)
			}))
			defer srv.Close()

			rp := skip_reverse_proxy.New(
				&url.URL{
					Scheme: "http",
					Host:   srv.URL[7:],
					Path:   "/%s/unique_path",
				},
				skip_reverse_proxy.WithGetUUIDFromPath(func(path string) string {
					return path[11:15]
				}),
			)
			require.NotNil(t, rp)

			t.Run("FirstRequest", func(t *testing.T) {
				request, err := http.NewRequestWithContext(ctx, "GET", "/something/1234/otherthing", nil)
				require.NoError(t, err)
				require.NotNil(t, request)

				w := httptest.NewRecorder()
				rp.ServeHTTP(w, request)

				resp := w.Result()
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "/1234/unique_path", path)
			})

			t.Run("SecondRequest", func(t *testing.T) {
				request2, err := http.NewRequestWithContext(ctx, "GET", "/something/5678/otherthing", nil)
				require.NoError(t, err)
				require.NotNil(t, request2)

				w := httptest.NewRecorder()
				rp.ServeHTTP(w, request2)

				resp := w.Result()
				assert.Equal(t, 200, resp.StatusCode)
				assert.Equal(t, "/5678/unique_path", path)
			})
		})
	})
}
