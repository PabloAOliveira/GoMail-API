.PHONY: build run docker-build docker-run docker-compose-up docker-compose-down

# Variáveis
APP_NAME = email-service

# Comandos locais
build:
	go build -o $(APP_NAME) .

run: build
	./$(APP_NAME)

# Comandos Docker
docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p 2000:2000 --env-file .env $(APP_NAME)

# Comandos Docker Compose
docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f

# Limpar recursos
clean:
	rm -f $(APP_NAME)
	docker-compose down -v

# Reiniciar aplicação com Docker Compose
restart:
	docker-compose restart email-service