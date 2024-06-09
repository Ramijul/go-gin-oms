FROM golang:latest
WORKDIR /app/payments
COPY go.mod go.sum ./
RUN go mod download
COPY . .
