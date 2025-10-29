package rabbitmq

import (
	"encoding/json"
	"log"

	"hospital_management_system/internal/infra/repository"
	"hospital_management_system/internal/models"
	"hospital_management_system/internal/pkg/helpers"

	"github.com/streadway/amqp"
)

func StartConsumer(amqpURL, queueName string, smtpHost string, smtpPort int, smtpUser, smtpPass string, emailRepo repository.EmailRepository) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Declare the queue before consuming
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var job helpers.EmailJob
			if err := json.Unmarshal(d.Body, &job); err != nil {
				log.Println("Failed to decode job:", err)
				continue
			}
			if err := helpers.SendEmail(job, smtpHost, smtpPort, smtpUser, smtpPass); err != nil {
				errMsg := err.Error()
				if err := emailRepo.UpdateEmailStatus(job.EmailID, models.EmailStatusFailed, &errMsg); err != nil {
					log.Println("Failed to update email status:", err)
				}
			} else {
				if err := emailRepo.UpdateEmailStatus(job.EmailID, models.EmailStatusSent, nil); err != nil {
					log.Println("Failed to update email status:", err)
				}
			}
		}
	}()

	log.Println("Email worker running...")
	<-forever
}
