FROM alpine:latest

# Instala certificados necesarios para conexiones HTTPS
RUN apk --no-cache add ca-certificates

# Limita el uso de memoria del heap a 500MB
ENV GOMEMLIMIT=500MB

# Establece el directorio de trabajo para el contenedor
WORKDIR /root/

# Copia el binario compilado (asegúrate de que 'main' esté en la misma carpeta que este Dockerfile)
COPY main .

# Expone el puerto en el que la aplicación escucha
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]
