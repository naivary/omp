package main

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	clubv1 "github.com/naivary/omp/api/club/v1"
)

func getAllClubProfiles() ([]*clubv1.Profile, error) {
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return nil, err
	}
	endpoint, err := url.JoinPath(baseURL, "clubs")
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	clubs := make([]*clubv1.Profile, 0, 100)
	return clubs, json.NewDecoder(res.Body).Decode(&clubs)
}

func randClubProfile(clubs []*clubv1.Profile) *clubv1.Profile {
	index := rand.Int64N(int64(len(clubs)))
	return clubs[index]
}

func TestCreateClub(t *testing.T) {
	tests := []struct {
		name string
		code int
		req  clubv1.CreateClubRequest
	}{
		{
			name: "valid request",
			code: http.StatusCreated,
			req: clubv1.CreateClubRequest{
				Email:    "info@omp.de",
				Password: "testpassword",
				Name:     "randomname",
				Location: "Berlin",
				Timezone: "Berlin/Europe",
			},
		},
	}
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}
	endpoint, err := url.JoinPath(baseURL, "clubs")
	if err != nil {
		t.Fatalf("URL join path: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest(http.MethodPost, endpoint, &tc.req)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}
		})
	}
}

func TestReadClub(t *testing.T) {
	clubs, err := getAllClubProfiles()
	if err != nil {
		t.Fatalf("get all club profiles: %v", err)
	}

	tests := []struct {
		name string
		code int
		club *clubv1.Profile
	}{
		{
			name: "valid request",
			code: http.StatusOK,
			club: randClubProfile(clubs),
		},
	}
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := strconv.FormatInt(tc.club.ID, 10)
			endpoint, err := url.JoinPath(baseURL, "clubs", id)
			if err != nil {
				t.Fatalf("URL join path: %s", err)
			}
			r := NewRequest[any](http.MethodGet, endpoint, nil)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}
		})
	}
}

func TestReadAllClubs(t *testing.T) {
	tests := []struct {
		name string
		code int
	}{
		{
			name: "valid request",
			code: http.StatusOK,
		},
	}
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}
	endpoint, err := url.JoinPath(baseURL, "clubs")
	if err != nil {
		t.Fatalf("URL join path: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRequest[any](http.MethodGet, endpoint, nil)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}
			// there will be 100 club profiles provisioned beforehand
			numClubs := 100
			clubs := make([]*clubv1.Profile, 0, numClubs)
			err = json.NewDecoder(res.Body).Decode(&clubs)
			if err != nil {
				t.Fatalf("json decode: %v", err)
			}
			if len(clubs) < 100 {
				t.Errorf("number of club profiles not expected. Got: %d. Want: >=100", len(clubs))
			}
		})
	}
}

func TestUpdateClub(t *testing.T) {
	clubs, err := getAllClubProfiles()
	if err != nil {
		t.Fatalf("get all club profiles: %v", err)
	}

	tests := []struct {
		name    string
		code    int
		club    *clubv1.Profile
		req     clubv1.UpdateClubRequest
		isValid func(before, after *clubv1.Profile) bool
	}{
		{
			name: "valid request",
			code: http.StatusNoContent,
			club: randClubProfile(clubs),
			req: clubv1.UpdateClubRequest{
				Name: "updated_name",
			},
			isValid: func(before, after *clubv1.Profile) bool {
				return before.Name != after.Name
			},
		},
	}
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := strconv.FormatInt(tc.club.ID, 10)
			endpoint, err := url.JoinPath(baseURL, "clubs", id)
			if err != nil {
				t.Fatalf("URL join path: %s", err)
			}
			r := NewRequest[any](http.MethodPatch, endpoint, &tc.req)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				t.Errorf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}

			// get the club profile and check if the updated content is
			// corrected
			getReq := NewRequest[any](http.MethodGet, endpoint, nil)
			updatedRes, err := http.DefaultClient.Do(getReq)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			updatedProfile := &clubv1.Profile{}
			err = json.NewDecoder(updatedRes.Body).Decode(updatedProfile)
			if err != nil {
				t.Fatalf("json decode: %v", err)
			}
			if !tc.isValid(tc.club, updatedProfile) {
				t.Errorf("updating of club profile did not work: %v", updatedProfile)
			}
		})
	}
}

func TestDeleteClub(t *testing.T) {
	clubs, err := getAllClubProfiles()
	if err != nil {
		t.Fatalf("get all club profiles: %v", err)
	}

	tests := []struct {
		name string
		code int
	}{
		{
			name: "valid request",
			code: http.StatusNoContent,
		},
	}
	ctx := context.Background()
	baseURL, err := NewTestServer(ctx, nil, os.Getenv, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf("new test sever: %s", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := strconv.FormatInt(tc.id, 10)
			endpoint, err := url.JoinPath(baseURL, "clubs", id)
			if err != nil {
				t.Fatalf("URL join path: %s", err)
			}
			r := NewRequest[any](http.MethodDelete, endpoint, nil)
			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.code {
				t.Fatalf("unexpected status code. Got: %d. Wanted: %d", res.StatusCode, tc.code)
			}

			// get the club profile and check if the updated content is
			// corrected
			getReq := NewRequest[any](http.MethodGet, endpoint, nil)
			updatedRes, err := http.DefaultClient.Do(getReq)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			if updatedRes.StatusCode != http.StatusBadRequest {
				t.Fatalf("unexpected status code. Got: %d. Want: %d", updatedRes.StatusCode, http.StatusNotFound)
			}
		})
	}
}
