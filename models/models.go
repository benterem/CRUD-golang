package models

type Container struct {
	ID          uint64 `json:"id"`
	ContainerId string `json:"containerId"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Size        string `json:"size"`
	Deleted     bool   `json:"-"`
}

type Response struct {
	Message string `json:"message,omitempty"`
}
