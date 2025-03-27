# Imagen base de Go
FROM golang:1.21

# Directorio de trabajo dentro del contenedor
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copiar todos los archivos al contenedor
COPY . .

# Compilar el programa
RUN go build -o app .

# Exponer el puerto 8080
EXPOSE 8080

# Comando para ejecutar el programa
CMD ["./app"]
