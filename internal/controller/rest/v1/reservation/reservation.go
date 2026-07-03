package reservation

import (
	reservationSrv "reservation/internal/service/reservation"

	"github.com/gofiber/fiber/v3"
)

type ReservationController struct {
	reservationSrv reservationSrv.ReservationService
}

func NewReservationController(reservationSrv reservationSrv.ReservationService) *ReservationController {
	return &ReservationController{reservationSrv: reservationSrv}
}

func (c *ReservationController) RegisterRoutes(router fiber.Router) {
	router.Get("/", c.GetReservations)
	router.Post("/", c.CreateReservation)
	router.Post("/:reservation_id/cancel", c.CancelReservation)
	router.Post("/:reservation_id/confirm", c.ConfirmReservation)
}
