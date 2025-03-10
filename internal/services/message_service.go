package services

import (
	"context"
	"errors"
	"mail-service/internal/models"
	"mail-service/internal/repository"
)

type MessageService struct {
	msgRepo  *repository.MessageRepository
	userRepo *repository.UserRepository
}

func NewMessageService(msgRepo *repository.MessageRepository, userRepo *repository.UserRepository) *MessageService {
	return &MessageService{
		msgRepo:  msgRepo,
		userRepo: userRepo,
	}
}

func (s *MessageService) SendMessage(senderID uint, receiverEmail, subject, body string) error {
	receiver, err := s.userRepo.GetUserByEmail(context.Background(), receiverEmail)
	if err != nil || receiver == nil {
		return errors.New("receiver not found")
	}

	msg := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiver.ID,
		Subject:    subject,
		Body:       body,
	}

	return s.msgRepo.CreateMessage(msg)
}

func (s *MessageService) GetMessages(userID uint) (sent []models.Message, received []models.Message, err error) {
	sent, err = s.msgRepo.GetSentMessages(userID)
	if err != nil {
		return
	}

	received, err = s.msgRepo.GetReceivedMessages(userID)
	return
}
