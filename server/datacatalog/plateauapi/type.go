package plateauapi

import (
	"strconv"
	"strings"
)

func to[T Node](n Node, err error) (t T, _ error) {
	if err != nil {
		return t, err
	}
	return n.(T), nil
}

type ID string

func NewID(id string, ty string) ID {
	return ID(id + ":" + ty)
}

func (i ID) String() string {
	return string(i)
}

func (i ID) ID() string {
	_, t, _ := strings.Cut(string(i), ":")
	return t
}

func (i ID) Type() string {
	t, _, _ := strings.Cut(string(i), ":")
	return t
}

type AreaCode string

func (a AreaCode) String() string {
	return string(a)
}

func (a AreaCode) Int() int {
	i, _ := strconv.Atoi(string(a))
	return i
}

func (a AreaCode) PrefectureCode() string {
	if len(a) < 2 {
		return ""
	}
	return string(a)[0:2]
}

func (a AreaCode) PrefectureCodeInt() int {
	i, _ := strconv.Atoi(a.PrefectureCode())
	return i
}

func (a AreaCode) IsPrefectureCode() bool {
	return len(a) == 2
}
