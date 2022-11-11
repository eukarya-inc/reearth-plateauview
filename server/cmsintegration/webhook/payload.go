package webhook

type Payload struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
