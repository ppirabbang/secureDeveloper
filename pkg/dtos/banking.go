package dtos

type DepositRequest struct {
	Amount int64 `json:"amount"`
}

type BalanceWithdrawRequest struct {
	Amount int64 `json:"amount"`
}

type TransferRequest struct {
	ToUsername string `json:"to_username"`
	Amount     int64  `json:"amount"`
}
