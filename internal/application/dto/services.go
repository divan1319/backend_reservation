package dto

type Service struct {
	Code          string `json:"code,omitempty"`
	Name          string `json:"name,omitempty"`
	EstimatedTime uint   `json:"estimated_time,omitempty"`
	Status        bool   `json:"status,omitempty"`
}
