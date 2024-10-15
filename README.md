# ToDoListWebApplication
![GitHub Created At](https://img.shields.io/github/created-at/franlo42/ToDoListWebApplication%20?color=%23A594F9)
![GitHub contributors](https://img.shields.io/github/contributors/franlo42/ToDoListWebApplication?COLOR=%23F95454)
![GitHub top language](https://img.shields.io/github/languages/top/franlo42/ToDoListWebApplication?color=%2377CDFF)
![Last commit](https://img.shields.io/github/last-commit/franlo42/ToDoListWebApplication?color=%2372BF78)
![GitHub repo size](https://img.shields.io/github/repo-size/franlo42/ToDoListWebApplication?color=%23FFBF61)

## Objective: 
Create a basic To-Do list web application using Go, Dockerize it, and set up multi-container orchestration with Docker Compose.
## Requirements:
1. **Go Application**:
  - Develop a basic REST API with Go that supports CRUD (Create, Read, Update, Delete) operations for managing to-do items.
  - Use the Gin framework for HTTP request routing.
2. **Docker**:
  - Write a `Dockerfile` to containerize the Go application.
  - Use a lightweight base image like `golang:alpine` for building the container.
  - Use a multi-stage build.
3. **Docker Compose**:
  - Create a `docker-compose.yml` file to define and run a multi-container setup.
  - The setup should include:
    - **Go App Container**: Containerize the Go web server.
    - **Database Container**: a database like PostgreSQL/MySQL, add a container for it.
  - Expose the necessary ports to access the web application.
4. **Extra Features** (Optional):
  - Use a volume to ensure data persistence for the database container.
  - Implement a logging mechanism using a Go package to log user interactions.
