package models

// GasInfo описывает данные о цене газа, оценке времени подтверждения и статусе сети
type GasInfo struct {
	Network       string  `json:"network"`
	GasPrice      float64 `json:"gas_price"`      // в Gwei
	EstimatedTime int     `json:"estimated_time"` // в секундах
	NetworkStatus string  `json:"network_status"` // например, "normal", "congested"
}
