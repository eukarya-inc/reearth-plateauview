package cmsintegration

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type ID struct {
	ItemID  string
	AssetID string
}

var ErrInvalidID = errors.New("invalid id")

func ParseID(id, secret string) (ID, error) {
	sig, payload, found := strings.Cut(id, ":")
	if !found {
		return ID{}, ErrInvalidID
	}

	itemID, assetID, found := strings.Cut(payload, ":")
	if !found {
		return ID{}, ErrInvalidID
	}

	if sig != sign(payload, secret) {
		return ID{}, ErrInvalidID
	}

	return ID{
		ItemID:  itemID,
		AssetID: assetID,
	}, nil
}

func (i ID) String(secret string) string {
	payload := fmt.Sprintf("%s:%s", i.ItemID, i.AssetID)
	sig := sign(payload, secret)
	return fmt.Sprintf("%s:%s", sig, payload)
}

func sign(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
