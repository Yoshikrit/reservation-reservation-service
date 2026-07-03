package product

import (
	"github.com/gofiber/fiber/v3"
)

func (c *ProductController) GetProduct(ctx fiber.Ctx) error {
	productID := ctx.Params("product_id")

	result, err := c.reservationSrv.GetProduct(ctx.Context(), productID)
	if err != nil {
		return err
	}

	return ctx.JSON(GetProductResponse{
		ProductID:      result.ProductID,
		Name:           result.Name,
		BasePrice:      result.BasePrice,
		StockTotal:     result.StockTotal,
		StockReserved:  result.StockReserved,
		StockAvailable: result.StockAvailable,
	})
}
