package dto

type Service struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	EstimatedTime uint   `json:"estimated_time"`
	Status        bool   `json:"status,omitempty"`
}
