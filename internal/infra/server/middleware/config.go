package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ConfigRequest é um middleware personalizado para logar requisições.
func ConfigRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	}
}
