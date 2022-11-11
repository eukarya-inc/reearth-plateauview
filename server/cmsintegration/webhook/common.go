package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/reearth/reearthx/util"
)

const (
	version = "v1"
	expires = time.Hour
	header  = "Reearth-Signature"
)

var ctxKey = struct{}{}

func AttacPayload(ctx context.Context, p *Payload) context.Context {
	return context.WithValue(ctx, ctxKey, p)
}

func GetPayload(ctx context.Context) *Payload {
	if c, ok := ctx.Value(ctxKey).(*Payload); ok {
		return c
	}
	return nil
}

func sign(payload, secret []byte, t time.Time, v string) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(fmt.Sprintf("%s:%d:", v, t.Unix())))
	_, _ = mac.Write(payload)
	s := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s,t=%d,%s", v, t.Unix(), s)
}

func validateSignature(actualSig string, payload, secret []byte) bool {
	if actualSig == "" {
		return false
	}

	sig := strings.Split(actualSig, ",")
	if len(sig) != 3 || sig[0] != version {
		return false
	}

	t, err := strconv.ParseInt(strings.TrimPrefix(sig[1], "t="), 10, 64)
	if err != nil {
		return false
	}

	timestamp := time.Unix(t, 0)

	if util.Now().Sub(timestamp) > expires {
		return false
	}

	expectedSig := sign(payload, secret, timestamp, version)
	return actualSig == expectedSig
}
