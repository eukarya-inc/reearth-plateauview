package cmswebhook

import (
	"encoding/json"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
)

type Payload struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	AssetData *AssetData      `json:"-"`
	ItemData  *ItemData       `json:"-"`
	Operator  Operator        `json:"operator"`
}

func (p *Payload) UnmarshalJSON(data []byte) error {
	type payload2 Payload
	if err := json.Unmarshal(data, (*payload2)(p)); err != nil {
		return err
	}
	if strings.HasPrefix(p.Type, "asset.") && p.Data != nil {
		p.AssetData = &AssetData{}
		if err := json.Unmarshal(p.Data, p.AssetData); err != nil {
			return err
		}
	} else if strings.HasPrefix(p.Type, "item.") && p.Data != nil {
		p.ItemData = &ItemData{}
		if err := json.Unmarshal(p.Data, p.ItemData); err != nil {
			return err
		}
	}
	p.Data = nil
	return nil
}

func (p Payload) ProjectID() string {
	if p.AssetData != nil {
		return p.AssetData.ProjectID
	}
	if p.ItemData != nil && p.ItemData.Schema != nil {
		return p.ItemData.Schema.ProjectID
	}
	return ""
}

type Operator struct {
	User        *User        `json:"user,omitempty"`
	Integration *Integration `json:"integration,omitempty"`
	Machine     *Machine     `json:"machine,omitempty"`
}

func (o Operator) IsUser() bool {
	return o.User != nil
}

func (o Operator) IsIntegration() bool {
	return o.Integration != nil
}

type User struct {
	ID string `json:"id"`
}

type Integration struct {
	ID string `json:"id"`
}

type Machine struct{}

type ItemData struct {
	Item   *cms.Item   `json:"item,omitempty"`
	Model  *cms.Model  `json:"model,omitempty"`
	Schema *cms.Schema `json:"schema,omitempty"`
}

type AssetData cms.Asset
