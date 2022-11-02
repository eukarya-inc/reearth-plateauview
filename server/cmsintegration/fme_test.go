package cmsintegration

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP_OK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	called := mockFMEServer(false)
	req := WebhookRequest("/api/webhook", map[string]any{}, "SECRET")
	w := httptest.NewRecorder()

	(&FMEHandler{}).ServeHTTP(w, req)
	assert.Equal(t, 200, w.Result().StatusCode)
	b := lo.Must(io.ReadAll(w.Result().Body))
	assert.Equal(t, `{"ok":true}`, string(b))
	assert.Equal(t, 1, called)
}

func TestFMEHandler_ServeHTTP_InvalidSecret(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	called := mockFMEServer(false)
	req := WebhookRequest("/api/webhook", map[string]any{}, "SECRET!")
	w := httptest.NewRecorder()

	(&FMEHandler{}).ServeHTTP(w, req)
	assert.Equal(t, 401, w.Result().StatusCode)
	b := lo.Must(io.ReadAll(w.Result().Body))
	assert.Equal(t, `{"error":"unauthorized"}`, string(b))
	assert.Equal(t, 0, called)
}
