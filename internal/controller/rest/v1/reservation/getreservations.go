package reservation

import (
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	reservationSrv "github.com/Yoshikrit/reservation/internal/service/reservation"

	"github.com/gofiber/fiber/v3"
)

func (c *ReservationController) GetReservations(ctx fiber.Ctx) error {
	var q GetReservationsQuery
	if err := ctx.Bind().Query(&q); err != nil {
		return apperror.NewError(40000000, err)
	}

	result, appErr := c.reservationSrv.GetReservations(ctx.Context(), &reservationSrv.ListReservationBody{
		ProductID: q.ProductID,
		Status:    q.Status,
		Limit:     q.Limit,
		Offset:    q.Offset,
	})
	if appErr != nil {
		return appErr
	}

	return ctx.JSON(result)
}
