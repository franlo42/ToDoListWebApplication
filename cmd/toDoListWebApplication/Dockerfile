# Etapa de construcción
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git && \
    go mod download && \
    go build -o todo-app && \
    apk del git

# Etapa de producción
FROM scratch
WORKDIR /
COPY --from=builder /app/todo-app /
EXPOSE 8080
CMD ["/todo-app"]
