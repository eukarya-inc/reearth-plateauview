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
	// TODO: ここで templates を Root に変換
	templateArray := []Template{}
	componentArray := []Component{}
	for i, t := range templates {
		templateArray = append(templateArray,
			Template{
				ID:       string(i),
				template: t.Fields[i].Value,
			})
	}
	for i, c := range data {
		componentArray = append(componentArray,
			Component{
				ID:        string(i),
				Component: c.Fields[i].Value,
			})
	}
	return Root{
		Templates:  templateArray,
		Components: componentArray,
	}
}
