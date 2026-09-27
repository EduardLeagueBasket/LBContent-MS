# Imagen base con Go 1.26
FROM golang:1.26-alpine

# Instalamos git y air para Hot Reload
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

WORKDIR /usr/src/app

# Copiamos archivos de dependencias de Go
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

EXPOSE 8080

CMD ["air", "--build.cmd", "go build -o ./tmp/main ./cmd/api", "--build.bin", "./tmp/main"]