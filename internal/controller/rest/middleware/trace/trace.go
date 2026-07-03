package trace

import (
	"context"
	"fmt"

	"reservation/internal/entity"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func Trace() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		traceID := ctx.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		event := fmt.Sprintf("Entry=API : %s %s", ctx.Method(), ctx.Path())
		reqCtx := context.WithValue(ctx.Context(), entity.ContextKeyEvent, event)
		reqCtx = context.WithValue(reqCtx, entity.ContextKeyTraceID, traceID)
		ctx.SetContext(reqCtx)
		ctx.Set("X-Request-ID", traceID)
		return ctx.Next()
	}
}
