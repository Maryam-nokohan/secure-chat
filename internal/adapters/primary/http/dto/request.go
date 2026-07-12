package dto

type SendMessageWSRequest struct {
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

type RegisterRequest struct{
	Username string `form:"username"`
	Password string `form:"password"`
	PublicKey string `form:"public_key"`
}
type LoginRequest struct{
	Username string `form:"username"`
	Password string `form:"password"`
}
type roomDTO struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Unread bool   `json:"unread"`
	}