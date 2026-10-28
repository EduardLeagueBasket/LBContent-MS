# Imagen base con Go 1.26
FROM golang:1.26-alpine

# Instalamos git y air para Hot Reload
RUN apk add --no-cache git \
    && git config --global --add safe.directory '*'
RUN go install github.com/air-verse/air@latest

# El repo se monta desde el host: git falla con "dubious ownership" (exit 128)
# y Go aborta el build al intentar estampar la revisión.
ENV GOFLAGS=-buildvcs=false

WORKDIR /usr/src/app

# Copiamos archivos de dependencias de Go
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

EXPOSE 8080

CMD ["air", "--build.cmd", "go build -buildvcs=false -o ./tmp/main ./cmd/api", "--build.bin", "./tmp/main"]