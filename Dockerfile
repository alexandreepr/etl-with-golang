FROM golang:1.22.2-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY . .

RUN go mod download

CMD ["air"]
