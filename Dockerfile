# Etapa de construcción
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh
RUN apk add --no-cache git && \
    go mod download && \
    go build -o todo-app ./cmd/toDoListWebApplication

# Etapa de producción
FROM scratch
WORKDIR /
COPY --from=builder /app/todo-app /
COPY --from=builder /app/wait-for-it.sh /wait-for-it.sh
EXPOSE 8080
CMD ["/wait-for-it.sh", "db:5432", "--", "/todo-app"]