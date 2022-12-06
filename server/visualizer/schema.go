package visualizer

import "github.com/eukarya-inc/reearth-plateauview/server/cms"

type Template struct {
	ID       string `json:"id"`
	template any    `json:"template`
}

type Component struct {
	ID        string `json:"id"`
	Component any    `json:"component"`
}

type Root struct {
	Templates  []Template  `json:"templates"`
	Components []Component `json::components"`
}

func ToTemplate(i cms.Item) Template {
	return Template{
		ID:       i.ID,
		template: i.Fields[0].Value,
	}
}

func ToComponent(i cms.Item) Component {
	return Component{
		ID:        i.ID,
		Component: i.Fields[0].Value,
	}
}

func ToRoot(templates []*cms.Item, data []*cms.Item) Root {
	// TODO: ここでRootに変換
	return Root{}
}
