package routes

import (
	"context"
	"strconv"
	"time"
	"todo-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SetupRoutes(app *fiber.App, dbPool *pgxpool.Pool) {
	app.Post("/tasks", func(c *fiber.Ctx) error {
		task := new(models.Task)

		if err := c.BodyParser(task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		query := `INSERT INTO tasks (title, description, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		err := dbPool.QueryRow(context.Background(), query, task.Title, task.Description, task.Status, time.Now(), time.Now()).Scan(&task.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to create task"})
		}

		return c.Status(fiber.StatusCreated).JSON(task)
	})

	app.Get("/tasks", func(c *fiber.Ctx) error {
		rows, err := dbPool.Query(context.Background(), "SELECT * FROM tasks")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to retrieve tasks"})
		}
		defer rows.Close()

		tasks := []models.Task{}
		for rows.Next() {
			var task models.Task
			err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to scan task"})
			}
			tasks = append(tasks, task)
		}

		return c.JSON(tasks)
	})

	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}

		task := new(models.Task)
		if err := c.BodyParser(task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		query := `UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5`
		_, err = dbPool.Exec(context.Background(), query, task.Title, task.Description, task.Status, time.Now(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to update task"})
		}

		return c.JSON(task)
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
		}

		query := `DELETE FROM tasks WHERE id = $1`
		_, err = dbPool.Exec(context.Background(), query, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to delete task"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
