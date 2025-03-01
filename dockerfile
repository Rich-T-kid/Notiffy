FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o notifService

EXPOSE 9999

CMD ["./notifService"]