package dto

type SendMessageWSRequest struct {
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

type RegisterRequest struct{
	Username string `form:"username"`
	Password string `form:"password"`
}
type LoginRequest struct{
	Username string `form:"username"`
	Password string `form:"password"`
}