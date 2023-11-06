package sdk

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFMEID(t *testing.T) {
	i := fmeID{ItemID: "item", AssetID: "asset", ProjectID: "project"}
	assert.Equal(t, i, lo.Must(parseFMEID(i.String("aaa"), "aaa")))
	_, err := parseFMEID(i.String("aaa"), "aaa2")
	assert.Same(t, ErrInvalidFMEID, err)
}
