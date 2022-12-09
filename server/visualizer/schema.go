package visualizer

import (
	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/samber/lo"
)

type Template struct {
	ID       string `json:"id"`
	Template any    `json:"template"`
}

type Data struct {
	ID        string `json:"id"`
	Component any    `json:"component"`
}

type Root struct {
	Templates  []Template `json:"templates"`
	Components []Data     `json:"data"`
}

func ToTemplate(i cms.Item) Template {
	return Template{
		ID:       i.ID,
		Template: i.Fields[0].Value,
	}
}

func ToComponent(i cms.Item) Data {
	return Data{
		ID:        i.ID,
		Component: i.Fields[0].Value,
	}
}

func ToRoot(t []*cms.Item, d []*cms.Item) *Root {
	templates := lo.Map(t, func(i *cms.Item, _ int) Template {
		return ToTemplate(*i)
	})

	components := lo.Map(d, func(i *cms.Item, _ int) Data {
		return ToComponent(*i)
	})
	return &Root{
		Components: components,
		Templates:  templates,
	}
}
