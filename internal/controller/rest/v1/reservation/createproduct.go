package reservation

import (
	"reservation/internal/pkg/apperror"

	reservationSrv "reservation/internal/service/reservation"

	"github.com/gofiber/fiber/v3"
)

func (c *ReservationController) CreateReservation(ctx fiber.Ctx) error {
	var request CreateReservationRequest
	if err := ctx.Bind().JSON(&request); err != nil {
		return apperror.NewError(40000000, err)
	}

	requestToSrv := parseCreateRequestToService(request)
	if err := c.reservationSrv.CreateReservation(ctx.Context(), &requestToSrv); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func parseCreateRequestToService(request CreateReservationRequest) reservationSrv.CreateReservationRequest {
	return reservationSrv.CreateReservationRequest{
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		TtlSecond: request.TtlSecond,
	}
}
