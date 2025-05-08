package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"email-service/models"
	"email-service/rabbitmq"
	"email-service/worker" 
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	// Iniciar o consumidor de mensagens em uma goroutine
	go worker.StartConsumer()

	router := gin.Default()

	router.POST("/send-email", func(c *gin.Context) {
		var task models.EmailTask

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		err := rabbitmq.Publish(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao publicar na fila"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Tarefa enfileirada com sucesso!"})
	})

	router.Run(":2000")
}