package models

// BalanceResponse описывает ответ с балансом кошелька
type BalanceResponse struct {
	Address    string `json:"address"`
	Blockchain string `json:"blockchain"`
	Balance    string `json:"balance"`
}
