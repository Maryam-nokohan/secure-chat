package message

import "github.com/maryam-nokohan/secure-chat/internal/core/ports"

type MessageService struct {
	messageRepo ports.MessageRepository
	userRepo    ports.UserRepository
}

func NewMessageService(
	msgRepo ports.MessageRepository,
	userRepo ports.UserRepository) *MessageService {
	return &MessageService{
		messageRepo: msgRepo,
		userRepo:    userRepo,
	}
}

