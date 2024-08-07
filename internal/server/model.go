package server

type APIRespond struct {
	Message string `json:"message,omitempty"`
}

type PublishRequest struct {
	Topic     string `query:"topic"`
	Partition int32  `query:"partition"`
	Key       string `query:"key"`
	Raw       bool   `query:"raw"`
}

type GroupRequest struct {
	GroupID string `json:"group_id"`
	Topic   string `json:"topic"`
	From    string `json:"from"`
}
