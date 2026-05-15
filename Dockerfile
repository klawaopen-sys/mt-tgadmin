FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o mt-tgadmin .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mt-tgadmin .
COPY .bot.yml .bot.yml

EXPOSE 8080
CMD ["./mt-tgadmin", "serve", "--settings", ".bot.yml"]
