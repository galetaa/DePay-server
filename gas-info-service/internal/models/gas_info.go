package models

// GasInfo описывает данные о цене газа, оценке времени подтверждения и статусе сети
type GasInfo struct {
	Network       string  `json:"network"`
	GasPrice      float64 `json:"gas_price"`      // в Gwei
	EstimatedTime int     `json:"estimated_time"` // в секундах
	NetworkStatus string  `json:"network_status"` // например, "normal", "congested"
}

type GasHistoryPoint struct {
	Network       string  `json:"network"`
	GasPrice      float64 `json:"gas_price"`
	EstimatedTime int     `json:"estimated_time"`
	NetworkStatus string  `json:"network_status"`
	CapturedAt    string  `json:"captured_at"`
}
