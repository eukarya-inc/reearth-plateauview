package cmsintegration

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
)

func WebhookRequest(target string, body any, secret string) *http.Request {
	b := bytes.NewReader(lo.Must(json.Marshal(body)))
	req := httptest.NewRequest("POST", target, b)
	// signature := TODO

	return req
}

func mockFMEServer(fail bool) func() int {
	count := 0

	httpmock.RegisterResponder("POST", "http://test-fme.example.com/fmejobsubmitter/test-repo/test-ws.fmw", func(req *http.Request) (*http.Response, error) {
		if err := parseFMEToken(req); err != nil {
			return nil, err
		}

		if fail {
			return httpmock.NewJsonResponse(400, map[string]any{
				"statusInfo": map[string]any{
					"message": "failure",
					"status":  "failure",
				},
			})
		}

		count++
		return httpmock.NewJsonResponse(200, map[string]any{
			"statusInfo": map[string]any{
				"message": "success",
				"status":  "success",
			},
		})
	})

	return func() int {
		return count
	}
}

func parseFMEToken(r *http.Request) error {
	aut := r.Header.Get("Authorization")
	_, token, found := strings.Cut(aut, "fmetoken token=")
	if !found || token != "TOKEN" {
		return errors.New("invalid fme token")
	}
	return nil
}
