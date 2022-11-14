package cmsintegration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFMEResult_GetResult(t *testing.T) {
	assert.Equal(t, "https://example.com", fmeResult{
		Results: map[string]string{
			"*": "https://example.com",
		},
	}.GetResult("*"))
	assert.Equal(t, "", fmeResult{
		Results: map[string]string{
			"**": "https://example.com",
		},
	}.GetResult("*"))
}

func TestNotifyHandler(t *testing.T) {
}
