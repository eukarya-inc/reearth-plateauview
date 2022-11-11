package webhook

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestEchoMiddleware(t *testing.T) {
	secret := []byte("abcdefg")
	e := echo.New()

	h := EchoMiddleware(secret)(func(c echo.Context) error {
		payload := GetPayload(c.Request().Context())
		return c.String(http.StatusOK, payload.Type+","+payload.Data.(string))
	})

	tests := []struct {
		name             string
		time             time.Time
		version          string
		payload          string
		wantCode         int
		wantBody         string
		invalidSignature bool
	}{
		{
			name:     "valid",
			time:     time.Now(),
			version:  "v1",
			payload:  `{"type":"asset.update","data":"aaaa"}`,
			wantCode: http.StatusOK,
			wantBody: "asset.update,aaaa",
		},
		{
			name:     "valid (old webhook)",
			time:     time.Now().Add(-time.Minute * 50),
			version:  "v1",
			payload:  `{"type":"asset.update","data":"aaaa"}`,
			wantCode: http.StatusOK,
			wantBody: "asset.update,aaaa",
		},
		{
			name:     "invalid timestamp",
			time:     time.Now().Add(-time.Minute * 60),
			version:  "v0",
			payload:  `{"type":"asset.update","data":"aaaa"}`,
			wantCode: http.StatusUnauthorized,
			wantBody: `"unauthorized"`,
		},
		{
			name:     "invalid version",
			time:     time.Now(),
			version:  "v0",
			payload:  `{"type":"asset.update","data":"aaaa"}`,
			wantCode: http.StatusUnauthorized,
			wantBody: `"unauthorized"`,
		},
		{
			name:     "invalid payload",
			time:     time.Now(),
			version:  "v1",
			payload:  "invalid",
			wantCode: http.StatusBadRequest,
			wantBody: `"invalid payload"`,
		},
		{
			name:             "invalid signature",
			time:             time.Now(),
			version:          "v1",
			invalidSignature: true,
			payload:          `{"type":"asset.update","data":"aaaa"}`,
			wantCode:         http.StatusUnauthorized,
			wantBody:         `"unauthorized"`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			payload := []byte(tc.payload)
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(payload))
			sig := Sign(payload, secret, tc.time, tc.version)
			if tc.invalidSignature {
				sig += ":::"
			}
			req.Header.Set("Reearth-Signature", sig)
			res := httptest.NewRecorder()

			assert.NoError(t, h(e.NewContext(req, res)))
			assert.Equal(t, tc.wantCode, res.Code)
			assert.Equal(t, tc.wantBody, strings.TrimSpace(string(lo.Must(io.ReadAll(res.Result().Body)))))
		})
	}
}
