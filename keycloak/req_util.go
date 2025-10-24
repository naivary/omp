package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func is4XX(code int) bool {
	return code >= 400 && code < 500
}

func newError(r io.Reader) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return err
	}
	return fmt.Errorf("response error: %s", buf.String())
}

func newRequest[T any](method, endpoint string, body T, header http.Header, query url.Values) (*http.Request, error) {
	if header == nil {
		header = http.Header{}
	}
	if query == nil {
		query = url.Values{}
	}
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(method, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	r.Header = header
	r.URL.RawQuery = query.Encode()
	return r, nil
}
