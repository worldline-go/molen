package server

type APIRespond struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}
