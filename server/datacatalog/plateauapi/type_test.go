package plateauapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	assert.Equal(t, ID("aaa:bbb"), NewID("bbb", "aaa"))
	assert.Equal(t, "aaa:bbb", ID("aaa:bbb").String())
	assert.Equal(t, Type("aaa"), ID("aaa:bbb").Type())
	assert.Equal(t, "bbb", ID("aaa:bbb").ID())
}

func TestAreaCode(t *testing.T) {
	assert.Equal(t, "01234567", AreaCode("01234567").String())

	assert.Equal(t, "01", AreaCode("01").PrefectureCode())
	assert.Equal(t, 1, AreaCode("01").PrefectureCodeInt())
	assert.True(t, AreaCode("01").IsPrefectureCode())

	assert.Equal(t, "01", AreaCode("01234567").PrefectureCode())
	assert.Equal(t, 1, AreaCode("01234567").PrefectureCodeInt())
	assert.False(t, AreaCode("01234567").IsPrefectureCode())

	assert.Equal(t, "", AreaCode("0").PrefectureCode())
	assert.Equal(t, 0, AreaCode("0").PrefectureCodeInt())
	assert.False(t, AreaCode("0").IsPrefectureCode())
}
