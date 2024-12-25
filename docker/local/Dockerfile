FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=1 go build -o app .
EXPOSE 8080

CMD go run main.go