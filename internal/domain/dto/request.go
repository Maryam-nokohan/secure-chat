package dto

type SendMessageWSRequest struct {
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

