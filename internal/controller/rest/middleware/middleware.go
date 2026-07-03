package middleware

import (
	"reservation/internal/controller/rest/middleware/cors"
	"reservation/internal/controller/rest/middleware/helmet"
	"reservation/internal/controller/rest/middleware/logger"
	"reservation/internal/controller/rest/middleware/recover"
	"reservation/internal/controller/rest/middleware/responsetime"
	"reservation/internal/controller/rest/middleware/trace"

	"github.com/gofiber/fiber/v3"
)

func NewMiddleware() []fiber.Handler {
	return []fiber.Handler{
		recover.Recover(),
		logger.Logger(),
		cors.Cors(),
		responsetime.ResponseTime(),
		helmet.Helmet(),
		trace.Trace(),
	}
}
