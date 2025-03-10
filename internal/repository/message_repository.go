package repository

import (
	"mail-service/internal/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateMessage(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *MessageRepository) GetSentMessages(senderID uint) ([]models.Message, error) {
	var messages []models.Message

	err := r.db.Where("sender_id = ?", senderID).Order("sent_at desc").Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) GetReceivedMessages(receiverID uint) ([]models.Message, error) {
	var messages []models.Message

	err := r.db.Where("receiver_id = ?", receiverID).Order("sent_at desc").Find(&messages).Error
	return messages, err
}
