package product

import (
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	reservationSrv "github.com/Yoshikrit/reservation/internal/service/reservation"

	"github.com/gofiber/fiber/v3"
)

func (c *ProductController) CreateProduct(ctx fiber.Ctx) error {
	var request CreateProductRequest
	if err := ctx.Bind().JSON(&request); err != nil {
		return apperror.NewError(40000000, err)
	}

	if err := c.reservationSrv.CreateProduct(ctx.Context(), &reservationSrv.CreateProductRequest{
		ProductID: request.ProductID,
		Name:      request.Name,
		Stock:     request.Stock,
		BasePrice: request.BasePrice,
	}); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{})
}
