package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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

	// Configuração do CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.POST("/send-email", func(c *gin.Context) {
		var task models.EmailTask

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		fmt.Println("Passou Aqui")
		
		err := rabbitmq.Publish(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao publicar na fila"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Tarefa enfileirada com sucesso!"})
	})

	router.Run(":2000")
}