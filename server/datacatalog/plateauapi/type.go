package plateauapi

import (
	"encoding/base64"
	"strconv"
	"strings"
)

type Type string

const (
	TypeDataset     Type = "dataset"
	TypeDatasetItem Type = "datasetitem"
	TypeDatasetType Type = "datasettype"
	TypeArea        Type = "area"
	TypePlateauSpec Type = "plateauspec"
)

func to[T Node](n Node, err error) (t T, _ error) {
	if err != nil {
		return t, err
	}
	u, ok := n.(T)
	if !ok {
		return t, nil
	}
	return u, nil
}

type ID string

const useBase64 = true

func NewID(id string, ty Type) ID {
	idstr := string(ty) + ":" + id
	if !useBase64 {
		return ID(idstr)
	}

	return ID(base64.StdEncoding.EncodeToString([]byte(idstr)))
}

func (i ID) String() string {
	return string(i)
}

func (i ID) ID() string {
	id, _ := i.Unwrap()
	return id
}

func (i ID) Type() Type {
	_, t := i.Unwrap()
	return t
}

func (i ID) Unwrap() (string, Type) {
	idstr := string(i)
	if useBase64 {
		di, err := base64.StdEncoding.DecodeString(string(i))
		if err != nil {
			return "", ""
		}

		idstr = string(di)
	}

	ty, id, _ := strings.Cut(idstr, ":")
	return id, Type(ty)
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
