# Etapa de construcción: Utiliza la imagen de Go 1.23.0 basada en Alpine para compilar la aplicación
FROM golang:1.23-alpine AS builder

# Instala herramientas necesarias (por ejemplo, git para descargar dependencias)
RUN apk update && apk add --no-cache git

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos de módulos y descarga las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código fuente al contenedor
COPY . .

# Compila la aplicación; se asume que el punto de entrada está en cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Etapa final: Imagen minimalista para correr la aplicación
FROM alpine:latest

# Instala certificados necesarios para conexiones HTTPS, si es requerido
RUN apk --no-cache add ca-certificates

# Establece el directorio de trabajo para el contenedor final
WORKDIR /root/

# Copia el binario compilado desde la etapa builder
COPY --from=builder /app/main .

# Expone el puerto en el que la aplicación escucha (asegúrate de que coincide con tu código)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
