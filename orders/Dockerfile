FROM golang:latest
WORKDIR /app/orders
COPY go.mod go.sum ./
RUN go mod download
COPY . .
