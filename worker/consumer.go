package worker

import (
	"encoding/json"
	"log"
	"os"

	"email-service/email"
	"email-service/models"

	"github.com/streadway/amqp"
)

func StartConsumer() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue", 
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %s", err)
	}

	// Consumir mensagens da fila
	msgs, err := ch.Consume(
		q.Name, // fila
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Erro ao consumir fila: %s", err)
	}

	log.Println("Consumidor iniciado. Esperando mensagens...")

	// Loop para processar as mensagens
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Mensagem recebida: %s", d.Body)
		
			var task models.EmailTask
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)
				continue
			}
			
			log.Printf("Enviando email para: %s, assunto: %s", task.To, task.Subject)
			err = email.SendEmail(task)
			if err != nil {
				log.Printf("Erro ao enviar e-mail: %s", err)
			} else {
				log.Printf("E-mail enviado com sucesso para %s", task.To)
			}
		}
	}()

	<-forever // Aguarda indefinidamente (mantém o consumidor em execução)
}