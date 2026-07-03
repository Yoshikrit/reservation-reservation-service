package product

import (
	reservationSrv "reservation/internal/service/reservation"

	"github.com/gofiber/fiber/v3"
)

type ProductController struct {
	reservationSrv reservationSrv.ReservationService
}

func NewProductController(svc reservationSrv.ReservationService) *ProductController {
	return &ProductController{reservationSrv: svc}
}

func (c *ProductController) RegisterRoutes(router fiber.Router) {
	router.Post("/", c.CreateProduct)
	router.Get("/:product_id", c.GetProduct)
}
