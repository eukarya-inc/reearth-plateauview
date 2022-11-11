package webhook

type Payload struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
