package reservation

import (
	"github.com/gofiber/fiber/v3"
)

func (c *ReservationController) ConfirmReservation(ctx fiber.Ctx) error {
	reservationID := ctx.Params("reservation_id")

	if err := c.reservationSrv.ConfirmReservation(ctx.Context(), reservationID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
}
