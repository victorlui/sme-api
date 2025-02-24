# Use a imagem oficial do Go como base
FROM golang:1.23

WORKDIR /app

COPY . .


COPY go.mod go.sum ./
RUN go mod tidy

RUN go build -v -o main ./cmd/main.go

# Exponha a porta em que o aplicativo vai rodar
EXPOSE 9090

# Comando para rodar o seu projeto Go
CMD ["/app/main"]
