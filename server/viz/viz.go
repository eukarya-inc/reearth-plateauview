package viz

import "net/http"

type Handler struct {
	// config
}

func NewHandler() *Handler {
	return &Handler{
		// config
	}
}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("implemented")
}
