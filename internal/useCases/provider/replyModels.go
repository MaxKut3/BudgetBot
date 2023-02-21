package provider

type ModelResponse struct {
	Success bool    `json:"success"`
	Date    string  `json:"date"`
	Result  float64 `json:"result"`
}
