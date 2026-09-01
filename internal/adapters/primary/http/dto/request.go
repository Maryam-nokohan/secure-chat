package dto

type SendMessageWSRequest struct {
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

type RegisterRequest struct {
	Username          string `form:"username"`
	Password          string `form:"password"`
	PublicKey         string `form:"public_key"`
	WrappedPrivateKey string `form:"wrapped_private_key"`
	PrivateKeyIV      string `form:"private_key_iv"`
	PrivateKeySalt    string `form:"private_key_salt"`
}
type LoginRequest struct{
	Username string `form:"username"`
	Password string `form:"password"`
}
