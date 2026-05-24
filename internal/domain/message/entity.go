package message

import (
	"time"

	"github.com/maryam-nokohan/secure-chat/internal/domain/user"
)

type Message struct {
	Sender  user.User
	Reciver user.User
	Text    string
	Date    time.Time
}
