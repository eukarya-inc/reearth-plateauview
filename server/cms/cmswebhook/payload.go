package cmswebhook

import "github.com/samber/lo"

type Payload struct {
	Type     string   `json:"type"`
	Data     Data     `json:"data"`
	Operator Operator `json:"operator"`
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

type Data struct {
	Item   *Item   `json:"item"`
	Schema *Schema `json:"schema"`
}

type Item struct {
	ID      string  `json:"id"`
	ModelID string  `json:"modelId"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type Schema struct {
	ID        string        `json:"id"`
	Fields    []SchemaField `json:"fields"`
	ProjectID string        `json:"projectId"`
}

type SchemaField struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Key  string `json:"key"`
}

func (d Item) Field(fid string) *Field {
	f, ok := lo.Find(d.Fields, func(f Field) bool {
		return f.ID == fid
	})
	if !ok {
		return nil
	}
	return &f
}

func (d Schema) FieldIDByKey(k string) string {
	f, ok := lo.Find(d.Fields, func(f SchemaField) bool {
		return f.Key == k
	})
	if !ok {
		return ""
	}
	return f.ID
}
