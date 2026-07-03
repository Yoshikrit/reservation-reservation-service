package product_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"reservation/config"
	"reservation/internal/controller/rest"
	ctrlRest "reservation/internal/controller/rest/v1/product"
	"reservation/internal/pkg/apperror"
	svc "reservation/internal/service/reservation"
	svcMocks "reservation/internal/service/reservation/mocks"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newApp(svcMock *svcMocks.ReservationService) *fiber.App {
	app := fiber.New(config.NewRestConfig(rest.ErrorHandler()))
	ctrl := ctrlRest.NewProductController(svcMock)
	ctrl.RegisterRoutes(app.Group("/products"))
	return app
}

// ── CreateProduct ──────────────────────────────────────────────────────────────

func TestREST_CreateProduct_Success(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("CreateProduct", mock.Anything, mock.MatchedBy(func(r *svc.CreateProductRequest) bool {
		return r.ProductID == "prod-001" && r.Name == "Laptop" && r.Stock == 10 && r.BasePrice == 999.99
	})).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"product_id": "prod-001",
		"name":       "Laptop",
		"stock":      10,
		"base_price": 999.99,
	})
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestREST_CreateProduct_ValidationError(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)

	body, _ := json.Marshal(map[string]any{"name": "Laptop"}) // missing product_id and base_price
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockSvc.AssertNotCalled(t, "CreateProduct")
}

func TestREST_CreateProduct_ServiceError(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("CreateProduct", mock.Anything, mock.Anything).
		Return(apperror.NewError(40900000, nil, "product", "prod-001"))

	body, _ := json.Marshal(map[string]any{
		"product_id": "prod-001",
		"name":       "Laptop",
		"stock":      10,
		"base_price": 999.99,
	})
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

// ── GetProduct ─────────────────────────────────────────────────────────────────

func TestREST_GetProduct_Success(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("GetProduct", mock.Anything, "prod-001").
		Return(&svc.GetProductResponse{
			ProductID:      "prod-001",
			Name:           "Laptop",
			BasePrice:      999.99,
			StockTotal:     10,
			StockReserved:  3,
			StockAvailable: 7,
		}, nil)

	req := httptest.NewRequest("GET", "/products/prod-001", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "prod-001", body["product_id"])
	assert.Equal(t, "Laptop", body["name"])
}

func TestREST_GetProduct_NotFound(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("GetProduct", mock.Anything, "no-such").
		Return((*svc.GetProductResponse)(nil), apperror.NewError(40400000, nil, "product", "no-such"))

	req := httptest.NewRequest("GET", "/products/no-such", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
