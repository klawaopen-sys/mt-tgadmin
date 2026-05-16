FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mt-tgadmin main.go

FROM alpine:latest

WORKDIR /app

# Копируем всё
COPY --from=builder /app/mt-tgadmin .
COPY .bot.yml .

# Проверяем, что файл есть
RUN ls -la

RUN chmod +x mt-tgadmin

EXPOSE 8080

CMD ["./mt-tgadmin", "run", "--settings", ".bot.yml"]
