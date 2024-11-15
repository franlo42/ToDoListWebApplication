package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func isValidStatus(status string) bool {
	return status == "pending" || status == "completed"
}

var db *sql.DB

func getTodos(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, status, created_at FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

func getTodosByStatus(c *gin.Context) {
	// Obtener el parámetro 'status' de los query params
	status := c.DefaultQuery("status", "") // Si no se pasa ningún valor, será una cadena vacía

	// Validar que el status sea válido
	if status != "" && !isValidStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido. Debe ser 'pending' o 'completed'"})
		return
	}

	// Consulta SQL para filtrar por status
	var rows *sql.Rows
	var err error
	if status == "" {
		// Si no se proporciona un filtro de status, devolver todos los todos
		rows, err = db.Query("SELECT id, title, status, created_at FROM todos")
	} else {
		// Si hay un filtro de status, buscar solo los que coinciden
		rows, err = db.Query("SELECT id, title, status, created_at FROM todos WHERE status = $1", status)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, todos)
}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo

	err := db.QueryRow("SELECT id, title, status, created_at FROM todos WHERE id = $1", id).Scan(
		&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar que el status esté permitido
	if !isValidStatus(newTodo.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido. Debe ser 'pending' o 'completed'"})
		return
	}

	err := db.QueryRow(
		"INSERT INTO todos (title, status, created_at) VALUES ($1, $2, $3) RETURNING id",
		newTodo.Title, newTodo.Status, time.Now(),
	).Scan(&newTodo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTodo)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	var updatedTodo Todo
	if err := c.BindJSON(&updatedTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidStatus(updatedTodo.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido. Debe ser 'pending' o 'completed'"})
		return
	}

	_, err := db.Exec(
		"UPDATE todos SET title = $1, status = $2 WHERE id = $3",
		updatedTodo.Title, updatedTodo.Status, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo actualizado con éxito"})
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo eliminado con éxito"})
}

func main() {
	connStr := "postgres://user:password@db:5432/tododb?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión a la base de datos:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("No se puede conectar a la base de datos:", err)
	}

	fmt.Println("Conexión exitosa a PostgreSQL")

	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos/status", getTodosByStatus) // Nueva ruta para obtener todos filtrados por 'status'
	router.GET("/todos/:id", getTodo)
	router.POST("/todos", createTodo)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)

	router.Run(":8080") // Escuchar en el puerto 8080
}
