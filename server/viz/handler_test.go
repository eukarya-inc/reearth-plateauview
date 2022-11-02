package viz

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// /api/viz/aaaa
	req := httptest.NewRequest("GET", "/aaaa", nil)
	w := httptest.NewRecorder()

	NewHandler().ServeHTTP(w, req)

	assert.Equal(t, 200, w.Result().StatusCode)
	assert.Equal(t, `{"test":"ok"}`, string(lo.Must(io.ReadAll(w.Result().Body))))
}
