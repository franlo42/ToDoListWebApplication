<div align="center"><a name="readme-top"></a>
  
  <img height="300" alt="To-Do List banner image" src="https://github.com/user-attachments/assets/30a161b1-8813-4fab-a7c0-fcd64c9e3ae0">
  
# To-Do List Web Application
  
  ![GitHub Created At](https://img.shields.io/github/created-at/franlo42/ToDoListWebApplication%20?color=%234F1787)
  ![GitHub contributors](https://img.shields.io/github/contributors/franlo42/ToDoListWebApplication?COLOR=%23FF6500)
  ![GitHub top language](https://img.shields.io/github/languages/top/franlo42/ToDoListWebApplication?color=%231230AE)
  ![Last commit](https://img.shields.io/github/last-commit/franlo42/ToDoListWebApplication?color=%23005B41)
  ![GitHub repo size](https://img.shields.io/github/repo-size/franlo42/ToDoListWebApplication?color=%23704264)

Simple web application that allows you to organize your daily tasks effectively. You can add new tasks, mark them as completed or delete them when you no longer need them. Ideal for those who want to keep a clear and accesible control of their to-dos. This application is your personal assistant to keep you focused and organized.

Managing your tasks has never been so easy!
</div>

<details>
<summary><kbd>Table of Contents</kbd></summary>

#### ToC

- [Objective](#-objective)
- [Requirements](#-requirements)
- [Quick Setup](#-quick-setup)
- [API Test](#-api-test)
- [Stopping the Application](#-stopping-the-application)

</details>

## 🎯 Objective

Create a basic To-Do list web application using Go, Dockerize it, and set up multi-container orchestration with Docker Compose.

## 📋 Requirements

1. **🦫 Go Application**
   - 📝 Develop a basic REST API with Go that supports CRUD (Create, Read, Update, Delete) operations for managing to-do items.
   - 🌐 Use the Gin framework for HTTP request routing.

2. **🐳 Docker**
   - 📄 Write a `Dockerfile` to containerize the Go application.
   - 📦 Use a lightweight base image like `golang:alpine` for building the container.
   - 🔄 Use a multi-stage build.

3. **🐙 Docker Compose**
   - 📄 Create a `docker-compose.yml` file to define and run a multi-container setup.
   - The setup should include:
     - 🫙 **Go App Container**: Containerize the Go web server.
     - 🗄️ **Database Container**: Use a database like PostgreSQL/MySQL, add a container for it.
   - 🌐 Expose the necessary ports to access the web application.

4. **✨ Extra Features** (Optional)
   - 💾 Use a volume to ensure data persistence for the database container.
   - 📊 Implement a logging mechanism using a Go package to log user interactions.

## ⚡ Quick Setup

You can run the application on your system using Docker Compose after cloning this repository:

```shell
git clone https://github.com/franlo42/ToDoListWebApplication.git
cd ToDoListWebApplication
docker compose up --build
```

## 💉 API Test

We can easily test the web app API functionalities with curl

**🗒️ Obtain the full list of ToDos**

```bash
curl -X GET http://localhost:8080/todos
```

**⚠️ Obtain the list of ToDos pending/completed**

```bash
curl -X GET http://localhost:8080/todos/status?status=pending
```

**➕ Add a new ToDo**

```bash
curl -X POST http://localhost:8080/todos -H "Content-Type: application/json" -d '{"title": "New Task", "status": "pending"}'
```

> [!IMPORTANT]  
> The **status attribute** must be 'pending' or 'completed'.

**🔄 Update a ToDo by ID**

```bash
curl -X PUT http://localhost:8080/todos/1 -H "Content-Type: application/json" -d '{"title": "Updated Task", "status": "completed"}'
```

**⁉️ Check a ToDo by ID**

```bash
curl -X GET http://localhost:8080/todos/1
```

**🗑️ Delete a ToDo by ID**

```bash
curl -X DELETE http://localhost:8080/todos/1
```

## 🛑 Stopping the Application

**1️⃣ Stop containers without deleting data**

```bash
docker-compose down
```

**2️⃣ Restart containers with persistent data**

```bash
docker-compose up
```
All previous data will still be available


> [!TIP]
> If you want to remove all data and restart fresh:
> ```bash
> docker-compose down -v 
> docker-compose up --build
> ```
> This will delete the volumes and reinitialize the database.
