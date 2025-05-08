FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum primeiro para aproveitar o cache do Docker
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar o restante do código
COPY . .

# Gerar a build
RUN go build -o email-service .

# Imagem final
FROM alpine:latest

WORKDIR /app

# Copiar o binário compilado da etapa de compilação
COPY --from=builder /app/email-service .
COPY --from=builder /app/.env .env

EXPOSE 2000

# Comando para iniciar o serviço
CMD ["./email-service"]