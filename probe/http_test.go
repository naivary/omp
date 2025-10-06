package probe

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/1s", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	mux.Handle("/timeout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	return mux
}

func TestDoHTTP(t *testing.T) {
	server := httptest.NewServer(newMux())
	defer server.Close()
	tests := []struct {
		name    string
		target  string
		timeout time.Duration
		res     Result
	}{
		{
			name:    "success",
			target:  "/1s",
			timeout: 5 * time.Second,
			res:     Success,
		},
		{
			name:    "timeout",
			target:  "/timeout",
			timeout: 1 * time.Second,
			res:     Failed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := url.JoinPath(server.URL, tc.target)
			if err != nil {
				t.Errorf("URL join path: %s", err)
				t.FailNow()
			}
			req, err := http.NewRequest(http.MethodGet, u, nil)
			if err != nil {
				t.Errorf("new request: %s", err)
				t.FailNow()
			}
			res, _ := DoHTTPWithClient(req, server.Client(), tc.timeout)
			if res != tc.res {
				t.Errorf("probe result not expected. Got: %d; Want: %d", res, tc.res)
			}
		})
	}
}
