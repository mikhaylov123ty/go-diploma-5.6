package models

type OrderData struct {
	OrderID    string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

type UserData struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}
