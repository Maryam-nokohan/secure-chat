package dto

type SendMessageWSResponse struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
	Date       string `json:"date"`
}
type AuthResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Redirect string `json:"redirect"`
}
