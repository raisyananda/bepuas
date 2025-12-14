package config

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	err := c.Next()
	method := c.Method()
	path := c.Path()

	fmt.Printf("[%s] %s - %v\n", method, path, time.Since(start))

	return err
}
