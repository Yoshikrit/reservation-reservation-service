package middleware

import (
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/cors"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/helmet"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/logger"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/metrics"
	otelMiddleware "github.com/Yoshikrit/reservation/internal/controller/rest/middleware/otel"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/recover"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/responsetime"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware/trace"

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
		otelMiddleware.Trace(),
		metrics.Metrics(),
	}
}
