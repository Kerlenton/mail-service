package mail

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type MailService struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewMailService() (*MailService, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL not set")
	}

	var conn *amqp.Connection
	var err error
	const maxAttempts = 5
	for i := 1; i <= maxAttempts; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection failed (attempt %d): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after retries: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"mail_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &MailService{
		conn:  conn,
		ch:    ch,
		queue: q,
	}, nil
}

func (ms *MailService) Close() {
	ms.ch.Close()
	ms.conn.Close()
}

func FormatEmail(to, subject, body string) string {
	return fmt.Sprintf("To: %s\nSubject: %s\n\n%s", to, subject, body)
}

func (ms *MailService) SendEmail(to, subject, body string) error {
	msg := FormatEmail(to, subject, body)
	err := ms.ch.Publish(
		"",
		ms.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send email to %s with subject %s: %w", to, subject, err)
	}
	log.Printf("Email successfully sent to %s with subject %s", to, subject)
	return nil
}
