FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mt-tgadmin main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mt-tgadmin .

RUN chmod +x mt-tgadmin

EXPOSE 8080

CMD sh -c 'printf "production: true\nbase_url: \"https://mt-tgadmin-production.up.railway.app\"\nbot_token: \"%s\"\nbot_chat_id: %s\ngui_password: \"%s\"\nwebserver_port: %s\nwebserver_hostname: \"%s\"\nwebserver_cookie_secret: \"%s\"\n" "$BOT_TOKEN" "$CHAT_ID" "$PASSWORD" "$PORT" "${HOST:-0.0.0.0}" "$COOKIE_SECRET" > .bot.yml && ./mt-tgadmin run --settings .bot.yml'
