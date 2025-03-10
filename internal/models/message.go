package models

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   uint      `gorm:"not null"`
	ReceiverID uint      `gorm:"not null"`
	Subject    string    `gorm:"not null"`
	Body       string    `gorm:"not null"`
	SentAt     time.Time `gorm:"autoCreateTime"`
}
