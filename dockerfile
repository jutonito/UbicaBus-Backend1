FROM alpine:latest

# Instala certificados necesarios para conexiones HTTPS, si es requerido
RUN apk --no-cache add ca-certificates

# Establece el directorio de trabajo para el contenedor
WORKDIR /root/

# Copia el binario compilado (asegúrate de que 'main' se encuentre en la misma carpeta que este Dockerfile)
COPY main .

# Expone el puerto en el que la aplicación escucha (asegúrate de que coincide con tu configuración)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
