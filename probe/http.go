package probe

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"time"
)

// DoHTTPWithClient sends a probe HTTP request using the provided client and waits
// until the server responds successfully or the timeout expires.
//
// It retries automatically if the connection is refused (for example, when the
// target server is still starting up).
func DoHTTPWithClient(r *http.Request, client *http.Client, timeout time.Duration) (Result, error) {
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()
	req := r.Clone(ctx)
	for {
		res, err := client.Do(req)
		if errors.Is(err, syscall.ECONNREFUSED) && ctx.Err() == nil {
			time.Sleep(250 * time.Millisecond)
			continue
		}
		if err := ctx.Err(); err != nil {
			return Failed, err
		}
		res.Body.Close()
		if isSuccessful(res.StatusCode) {
			return Success, nil
		}
	}
}

// DoHTTP sends an HTTP probe request using the default HTTP client and returns the
// resulting Status and any error encountered.
//
// It is a convenience wrapper around DoHTTPWithClient, which allows customization
// of the HTTP client used for the request.
func DoHTTP(r *http.Request, timeout time.Duration) (Result, error) {
	return DoHTTPWithClient(r, http.DefaultClient, timeout)
}

// isSuccessful reports whether the given HTTP status code indicates success.
//
// A status code is considered successful if it falls within the 2xx range
// (i.e., from http.StatusOK to http.StatusIMUsed, inclusive).
func isSuccessful(code int) bool {
	return code >= http.StatusOK && code <= http.StatusIMUsed
}
