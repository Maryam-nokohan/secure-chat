package dto

type SendMessageWSRequest struct {
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

type RegisterRequest struct {
	Username          string `json:"username" form:"username"`
	Email             string `json:"email" form:"email"`
	Password          string `json:"password" form:"password"`
	PublicKey         string `json:"public_key" form:"public_key"`
	WrappedPrivateKey string `json:"wrapped_private_key" form:"wrapped_private_key"`
	PrivateKeyIV      string `json:"private_key_iv" form:"private_key_iv"`
	PrivateKeySalt    string `json:"private_key_salt" form:"private_key_salt"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
