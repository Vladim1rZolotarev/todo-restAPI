package main

import (
	"context"
	"log"
	"todo-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	app := fiber.New()

	// Подключение к БД
	connString := "postgres://postgres:2739333@localhost:5432/tasks"
	log.Println("Connecting to database with URL:", connString)

	dbPool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer dbPool.Close()

	// Инициализация маршрутов
	routes.SetupRoutes(app, dbPool)

	// Запуск сервера fiber на порту 3000
	log.Fatal(app.Listen(":3000"))
}
