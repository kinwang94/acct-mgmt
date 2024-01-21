package model

type RequestPayload struct {
	Username string `json:"username" binding:"required" extensions:"x-order=1"`
	Password string `json:"password" binding:"required" extensions:"x-order=2"`
}

type ResponsePayload struct {
	Success bool   `json:"success" binding:"required" extensions:"x-order=1"`
	Reason  string `json:"reason" binding:"required" extensions:"x-order=2"`
}
