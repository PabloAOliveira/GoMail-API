
package rabbitmq

import (
	"log"
	"os"
	"github.com/streadway/amqp"
	"email-service/models"
	"github.com/joho/godotenv"
)

func Publish(task models.EmailTask) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %s", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao criar canal: %s", err)
		return err
	}
	defer ch.Close()


	// Declarar a fila
	q, err := ch.QueueDeclare(
		"email_queue", // Nome da fila
		true,          // Durável (a fila sobreviverá a reinicializações do RabbitMQ)
		false,         // Exclusiva (não será removida quando o cliente se desconectar)
		false,         // Não espera que outros consumidores cheguem
		false,         // Não espera que outros consumidores escutem
		nil,           // Propriedades adicionais
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %s", err)
		return err
	}

	body := []byte(`{"to":"` + task.To + `","subject":"` + task.Subject + `","body":"` + task.Body + `"}`)

		// Publicar a mensagem na fila
		err = ch.Publish(
			"",           // Exchange (usando a fila diretamente)
			q.Name,       // Nome da fila
			false,        // Mandatory
			false,        // Immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Fatalf("Erro ao publicar na fila: %s", err)
			return err
		}
	
	log.Println("Tarefa enfileirada com sucesso!")
	return nil
}
