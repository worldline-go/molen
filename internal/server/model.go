package server

type APIRespond struct {
	Message string `json:"message,omitempty"`
}

type PublishRequest struct {
	Topic     string `json:"topic"     param:"topic"`
	Partition int32  `json:"partition" param:"partition"`
	Key       string `json:"key"       query:"key"`
	Raw       bool   `json:"raw"       query:"raw"`
}

type GroupRequest struct {
	GroupID string `json:"group_id"`
	Topic   string `json:"topic"`
	From    string `json:"from"`
}
