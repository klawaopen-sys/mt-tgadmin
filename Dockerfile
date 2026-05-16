FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .

RUN chmod +x mt-tgadmin
RUN go mod download

# Если бинарник уже есть — используем его
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/mt-tgadmin .
COPY .bot.yml .

RUN chmod +x mt-tgadmin

EXPOSE 8080

CMD ["./mt-tgadmin", "run"]
