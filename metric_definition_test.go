package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	metricv1 "github.com/naivary/omp/api/metric/v1"
)

func TestCreateMetricDefinition(t *testing.T) {
	tests := []struct {
		name    string
		payload metricv1.CreateMetricRequest
		code    int
	}{
		{
			name: "create valid definition",
			payload: metricv1.CreateMetricRequest{
				Name:        "global_team_metric",
				Type:        metricv1.TypeCounter,
				Scope:       metricv1.ScopeTeam,
				Description: "this is a global metric which can be used by any of the teams to track a certain factor.",
			},
			code: http.StatusCreated,
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Errorf("new test sever: %s", err)
		t.FailNow()
	}
	endpoint, err := url.JoinPath(baseURL, "metrics/definition")
	if err != nil {
		t.Errorf("URL join path: %s", err)
		t.FailNow()
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := NewRequest(http.MethodPost, endpoint, tc.payload)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("send request: %v", err)
				t.FailNow()
			}
			if res.StatusCode != tc.code {
				t.Fatalf("unexpected status code. Got: %d; Want: %d", res.StatusCode, tc.code)
			}
		})
	}
}
