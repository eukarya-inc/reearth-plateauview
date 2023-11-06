package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ fmeInterface = (*fme)(nil)

func TestSignFMEID(t *testing.T) {
	payload := "xxxxxxxxx"
	signed := signFMEID(payload, "aaa")
	unsigned, err := unsignFMEID(signed, "aaa")
	assert.NoError(t, err)
	assert.Equal(t, payload, unsigned)

	unsigned2, err := unsignFMEID(signed, "aaa2")
	assert.Empty(t, unsigned2)
	assert.Same(t, ErrInvalidFMEID, err)
}
