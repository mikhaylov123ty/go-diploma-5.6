package models

import "time"

type OrderData struct {
	OrderID    string    `json:"number"`
	UserLogin  string    `json:"user_login,omitempty"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserData struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type BalanceData struct {
	UserLogin string  `json:"user_login,omitempty"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawData struct {
	UserLogin   string    `json:"user_login,omitempty"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type AccrualOrders struct {
	OrderID string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
