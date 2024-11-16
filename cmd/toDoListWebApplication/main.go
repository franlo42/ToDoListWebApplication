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
	c.Writer.Write([]byte("\n"))
}

func getTodosByStatus(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	if status != "" && !isValidStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido. Debe ser 'pending' o 'completed'"})
		return
	}

	var rows *sql.Rows
	var err error
	if status == "" {
		rows, err = db.Query("SELECT id, title, status, created_at FROM todos")
	} else {
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
	c.Writer.Write([]byte("\n"))
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
	c.Writer.Write([]byte("\n"))
}

func createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	c.Writer.Write([]byte("\n"))
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
	c.Writer.Write([]byte("\n"))
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo eliminado con éxito"})
	c.Writer.Write([]byte("\n"))
}

// Función para esperar la conexión a la base de datos
func waitForDB(connStr string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("Error al intentar conectar a la base de datos:", err)
		} else {
			if err = db.Ping(); err == nil {
				break
			} else {
				log.Println("No se pudo conectar a la base de datos:", err)
			}
		}
		log.Println("Esperando 2 segundos antes de intentar nuevamente...")
		time.Sleep(2 * time.Second)
	}
	return db, nil
}

func printLogo() {
	logo := `
 _____    ______      _     _     _   _    _      _      ___              
|_   _|   |  _  \    | |   (_)   | | | |  | |    | |    / _ \             
  | | ___ | | | |___ | |    _ ___| |_| |  | | ___| |__ / /_\ \_ __  _ __  
  | |/ _ \| | | / _ \| |   | / __| __| |/\| |/ _ \ '_ \|  _  | '_ \| '_ \ 
  | | (_) | |/ / (_) | |___| \__ \ |_\  /\  /  __/ |_) | | | | |_) | |_) |
  \_/\___/|___/ \___/\_____/_|___/\__|\/  \/ \___|_.__/\_| |_/ .__/| .__/ 
                                                             | |   | |    
                                                             |_|   |_|    
`
	fmt.Println(logo)
}

func main() {
	connStr := "postgres://user:password@db:5432/tododb?sslmode=disable"

	// Intentar conectar a la base de datos
	var err error
	db, err = waitForDB(connStr)
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos después de varios intentos:", err)
	}
	defer db.Close()

	fmt.Println("Conexión exitosa a PostgreSQL")

	printLogo()

	router := gin.Default()

	router.GET("/todos", getTodos)
	router.GET("/todos", getTodosByStatus)
	router.GET("/todos/:id", getTodo)
	router.POST("/todos", createTodo)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)

	router.Run(":8080")
}
