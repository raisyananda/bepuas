package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	latency := time.Since(start)

	// print info dasar
	fmt.Printf("[%s] %s %s - %d (%s)\n", start.Format(time.RFC3339), c.Method(), c.Path(), c.Response().StatusCode(), latency)

	return err
}
