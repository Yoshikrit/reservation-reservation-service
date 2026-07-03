package reservation_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"reservation/config"
	"reservation/internal/controller/rest"
	ctrlRest "reservation/internal/controller/rest/v1/reservation"
	"reservation/internal/pkg/apperror"
	svc "reservation/internal/service/reservation"
	svcMocks "reservation/internal/service/reservation/mocks"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newApp(svcMock *svcMocks.ReservationService) *fiber.App {
	app := fiber.New(config.NewRestConfig(rest.ErrorHandler()))
	ctrl := ctrlRest.NewReservationController(svcMock)
	ctrl.RegisterRoutes(app.Group("/reservations"))
	return app
}

// ── CreateReservation ──────────────────────────────────────────────────────────

func TestREST_CreateReservation_Success(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("CreateReservation", mock.Anything, mock.MatchedBy(func(r *svc.CreateReservationRequest) bool {
		return r.ProductID == "prod-001" && r.Quantity == 5 && r.TtlSecond == 60
	})).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"product_id": "prod-001",
		"quantity":   5,
		"ttl_second": 60,
	})
	req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestREST_CreateReservation_ValidationError(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)

	body, _ := json.Marshal(map[string]any{"product_id": "prod-001"}) // missing quantity and ttl_second
	req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockSvc.AssertNotCalled(t, "CreateReservation")
}

func TestREST_CreateReservation_ServiceError(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("CreateReservation", mock.Anything, mock.Anything).
		Return(apperror.NewError(42200000, nil))

	body, _ := json.Marshal(map[string]any{
		"product_id": "prod-001",
		"quantity":   5,
		"ttl_second": 60,
	})
	req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := newApp(mockSvc).Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// ── GetReservations ────────────────────────────────────────────────────────────

func TestREST_GetReservations_Success(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	now := time.Now()
	mockSvc.On("GetReservations", mock.Anything, mock.Anything).
		Return(&svc.ListReservationResponse{
			Reservations: []svc.ReservationItem{
				{ReservationID: "rsv-001", ProductID: "prod-001", ExpiresAt: now},
			},
			Pagination: svc.PaginationInfo{Limit: 20, Count: 1},
		}, nil)

	req := httptest.NewRequest("GET", "/reservations", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	reservations := body["reservations"].([]any)
	assert.Len(t, reservations, 1)
}

func TestREST_GetReservations_ServiceError(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("GetReservations", mock.Anything, mock.Anything).
		Return((*svc.ListReservationResponse)(nil), apperror.NewError(50000000, nil))

	req := httptest.NewRequest("GET", "/reservations", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

// ── CancelReservation ──────────────────────────────────────────────────────────

func TestREST_CancelReservation_Success(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("CancelReservation", mock.Anything, "rsv-001").Return(nil)

	req := httptest.NewRequest("POST", "/reservations/rsv-001/cancel", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestREST_CancelReservation_NotFound(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("CancelReservation", mock.Anything, "rsv-x").
		Return(apperror.NewError(40400000, nil, "reservation", "rsv-x"))

	req := httptest.NewRequest("POST", "/reservations/rsv-x/cancel", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

// ── ConfirmReservation ─────────────────────────────────────────────────────────

func TestREST_ConfirmReservation_Success(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("ConfirmReservation", mock.Anything, "rsv-001").Return(nil)

	req := httptest.NewRequest("POST", "/reservations/rsv-001/confirm", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestREST_ConfirmReservation_Expired(t *testing.T) {
	mockSvc := new(svcMocks.ReservationService)
	mockSvc.On("ConfirmReservation", mock.Anything, "rsv-001").
		Return(apperror.NewError(42200000, nil))

	req := httptest.NewRequest("POST", "/reservations/rsv-001/confirm", nil)
	resp, err := newApp(mockSvc).Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}
