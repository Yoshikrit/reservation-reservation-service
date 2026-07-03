package reservation_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Yoshikrit/reservation/internal/entity"
	inventoryGateway "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory"
	gatewayMocks "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory/mocks"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	cfgMocks "github.com/Yoshikrit/reservation/internal/repository/cfgpromotion/mocks"
	outboxMocks "github.com/Yoshikrit/reservation/internal/repository/outbox/mocks"
	rsvRepoMocks "github.com/Yoshikrit/reservation/internal/repository/reservation/mocks"
	"github.com/Yoshikrit/reservation/internal/service/constant"
	svc "github.com/Yoshikrit/reservation/internal/service/reservation"
	svcMocks "github.com/Yoshikrit/reservation/internal/service/reservation/mocks"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newRedis(t *testing.T) *redis.Client {
	t.Helper()
	mr := miniredis.RunT(t)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func trPassthrough(trMgr *svcMocks.TrManager) {
	trMgr.On("Do", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error)
			fn(context.Background())
		}).Return(nil)
}

// ── CreateProduct ──────────────────────────────────────────────────────────────

func TestCreateProduct_Success(t *testing.T) {
	gateway := new(gatewayMocks.InventoryGateway)
	gateway.On("CreateProduct", mock.Anything, mock.MatchedBy(func(r inventoryGateway.CreateProductRequest) bool {
		return r.ProductID == "prod-001" && r.Stock == 10
	})).Return(nil)

	s := svc.NewReservationService(nil, nil, nil, gateway, newRedis(t), nil)
	err := s.CreateProduct(context.Background(), &svc.CreateProductRequest{
		ProductID: "prod-001", Name: "Laptop", Stock: 10, BasePrice: 999.99,
	})

	assert.Nil(t, err)
	gateway.AssertExpectations(t)
}

func TestCreateProduct_GatewayError(t *testing.T) {
	gateway := new(gatewayMocks.InventoryGateway)
	gateway.On("CreateProduct", mock.Anything, mock.Anything).
		Return(apperror.NewError(40900000, nil, "product", "prod-001"))

	s := svc.NewReservationService(nil, nil, nil, gateway, newRedis(t), nil)
	err := s.CreateProduct(context.Background(), &svc.CreateProductRequest{
		ProductID: "prod-001", Name: "Laptop", Stock: 10, BasePrice: 999.99,
	})

	assert.NotNil(t, err)
	assert.Equal(t, apperror.CategoryConflict, err.Category)
}

// ── GetProduct ─────────────────────────────────────────────────────────────────

func TestGetProduct_Success(t *testing.T) {
	gateway := new(gatewayMocks.InventoryGateway)
	rsvRepo := new(rsvRepoMocks.ReservationRepository)

	gateway.On("GetProductByID", mock.Anything, "prod-001").
		Return(&inventoryGateway.GetProductByIDResponse{
			ProductID: "prod-001", Name: "Laptop", BasePrice: 999.99, Stock: 10,
		}, nil)
	rsvRepo.On("SumReservedQuantity", mock.Anything, "prod-001").Return(uint(3), nil)

	s := svc.NewReservationService(rsvRepo, nil, nil, gateway, newRedis(t), nil)
	resp, err := s.GetProduct(context.Background(), "prod-001")

	assert.Nil(t, err)
	assert.Equal(t, "prod-001", resp.ProductID)
	assert.Equal(t, uint(10), resp.StockTotal)
	assert.Equal(t, uint(3), resp.StockReserved)
	assert.Equal(t, uint(7), resp.StockAvailable)
}

func TestGetProduct_NotFound(t *testing.T) {
	gateway := new(gatewayMocks.InventoryGateway)
	gateway.On("GetProductByID", mock.Anything, "no-such").
		Return((*inventoryGateway.GetProductByIDResponse)(nil), apperror.NewError(40400000, nil, "product", "no-such"))

	s := svc.NewReservationService(nil, nil, nil, gateway, newRedis(t), nil)
	resp, err := s.GetProduct(context.Background(), "no-such")

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, apperror.CategoryNotFound, err.Category)
}

// ── GetReservations ────────────────────────────────────────────────────────────

func TestGetReservations_Success(t *testing.T) {
	rsvRepo := new(rsvRepoMocks.ReservationRepository)
	now := time.Now()

	rsvRepo.On("Count", mock.Anything, mock.Anything).Return(int64(2), nil)
	rsvRepo.On("FilterWithStatusOrder", mock.Anything, mock.Anything, 20, 0, true).
		Return([]entity.Reservation{
			{ReservationID: "rsv-001", ProductID: "prod-001", Status: constant.StatusHeld, ExpiresAt: now},
			{ReservationID: "rsv-002", ProductID: "prod-001", Status: constant.StatusHeld, ExpiresAt: now},
		}, nil)

	s := svc.NewReservationService(rsvRepo, nil, nil, nil, newRedis(t), nil)
	resp, err := s.GetReservations(context.Background(), &svc.ListReservationBody{ProductID: "prod-001"})

	assert.Nil(t, err)
	assert.Len(t, resp.Reservations, 2)
	assert.Equal(t, int64(2), resp.Pagination.Count)
}

// ── CreateReservation ──────────────────────────────────────────────────────────

func TestCreateReservation_Success(t *testing.T) {
	rsvRepo := new(rsvRepoMocks.ReservationRepository)
	gateway := new(gatewayMocks.InventoryGateway)
	cfgRepo := new(cfgMocks.CfgPromotionDateRepository)

	rsvRepo.On("SumReservedQuantity", mock.Anything, "prod-001").Return(uint(0), nil)
	gateway.On("CheckAndHold", mock.Anything, mock.Anything).
		Return(&inventoryGateway.CheckAndHoldResponse{Available: true, Price: 999.99}, nil)
	cfgRepo.On("Filter", mock.Anything, mock.Anything, 0, 0, true).
		Return([]entity.CfgPromotionDate{}, nil)
	rsvRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	s := svc.NewReservationService(rsvRepo, cfgRepo, nil, gateway, newRedis(t), nil)
	err := s.CreateReservation(context.Background(), &svc.CreateReservationRequest{
		ProductID: "prod-001", Quantity: 5, TtlSecond: 60,
	})

	assert.Nil(t, err)
	rsvRepo.AssertExpectations(t)
	gateway.AssertExpectations(t)
}

func TestCreateReservation_InsufficientStock(t *testing.T) {
	rsvRepo := new(rsvRepoMocks.ReservationRepository)
	gateway := new(gatewayMocks.InventoryGateway)

	rsvRepo.On("SumReservedQuantity", mock.Anything, "prod-001").Return(uint(0), nil)
	gateway.On("CheckAndHold", mock.Anything, mock.Anything).
		Return(&inventoryGateway.CheckAndHoldResponse{Available: false}, nil)

	s := svc.NewReservationService(rsvRepo, nil, nil, gateway, newRedis(t), nil)
	err := s.CreateReservation(context.Background(), &svc.CreateReservationRequest{
		ProductID: "prod-001", Quantity: 999, TtlSecond: 60,
	})

	assert.NotNil(t, err)
	assert.Equal(t, apperror.CategoryUnprocessable, err.Category)
}

func TestCreateReservation_LockRedisError(t *testing.T) {
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	mr.Close()

	s := svc.NewReservationService(nil, nil, nil, nil, rdb, nil)
	err := s.CreateReservation(context.Background(), &svc.CreateReservationRequest{
		ProductID: "prod-001", Quantity: 5, TtlSecond: 60,
	})

	assert.NotNil(t, err)
	assert.Equal(t, apperror.CategoryInternal, err.Category)
}

// ── CancelReservation ──────────────────────────────────────────────────────────

func TestCancelReservation_Success(t *testing.T) {
	rsvRepo := new(rsvRepoMocks.ReservationRepository)
	trMgr := new(svcMocks.TrManager)
	trPassthrough(trMgr)

	rsvRepo.On("FindForUpdate", mock.Anything, "rsv-001").
		Return(&entity.Reservation{ReservationID: "rsv-001", Status: constant.StatusHeld}, nil)
	rsvRepo.On("Patch", mock.Anything, mock.MatchedBy(func(r *entity.Reservation) bool {
		return r.ReservationID == "rsv-001" && r.Status == constant.StatusCancelled
	})).Return(nil)

	s := svc.NewReservationService(rsvRepo, nil, nil, nil, newRedis(t), trMgr)
	err := s.CancelReservation(context.Background(), "rsv-001")

	assert.Nil(t, err)
	rsvRepo.AssertExpectations(t)
}

func TestCancelReservation_InvalidStatus(t *testing.T) {
	trMgr := new(svcMocks.TrManager)
	trMgr.On("Do", mock.Anything, mock.Anything).
		Return(apperror.NewError(42200000, errors.New("only HELD reservations can be cancelled")))

	s := svc.NewReservationService(nil, nil, nil, nil, newRedis(t), trMgr)
	err := s.CancelReservation(context.Background(), "rsv-001")

	assert.NotNil(t, err)
	assert.Equal(t, apperror.CategoryUnprocessable, err.Category)
}

// ── ConfirmReservation ─────────────────────────────────────────────────────────

func TestConfirmReservation_Success(t *testing.T) {
	rsvRepo := new(rsvRepoMocks.ReservationRepository)
	outbox := new(outboxMocks.OutboxRepository)
	trMgr := new(svcMocks.TrManager)
	trPassthrough(trMgr)

	rsvRepo.On("FindForUpdate", mock.Anything, "rsv-001").
		Return(&entity.Reservation{
			ReservationID: "rsv-001",
			ProductID:     "prod-001",
			Status:        constant.StatusHeld,
			ExpiresAt:     time.Now().Add(10 * time.Minute),
		}, nil)
	rsvRepo.On("Patch", mock.Anything, mock.MatchedBy(func(r *entity.Reservation) bool {
		return r.ReservationID == "rsv-001" && r.Status == constant.StatusConfirmed
	})).Return(nil)
	outbox.On("Create", mock.Anything, mock.Anything).Return(nil)

	s := svc.NewReservationService(rsvRepo, nil, outbox, nil, newRedis(t), trMgr)
	err := s.ConfirmReservation(context.Background(), "rsv-001")

	assert.Nil(t, err)
	rsvRepo.AssertExpectations(t)
	outbox.AssertExpectations(t)
}

func TestConfirmReservation_Expired(t *testing.T) {
	rsvRepo := new(rsvRepoMocks.ReservationRepository)
	trMgr := new(svcMocks.TrManager)
	trPassthrough(trMgr)

	rsvRepo.On("FindForUpdate", mock.Anything, "rsv-001").
		Return(&entity.Reservation{
			ReservationID: "rsv-001",
			Status:        constant.StatusHeld,
			ExpiresAt:     time.Now().Add(-1 * time.Minute),
		}, nil)
	rsvRepo.On("Patch", mock.Anything, mock.MatchedBy(func(r *entity.Reservation) bool {
		return r.ReservationID == "rsv-001" && r.Status == constant.StatusCancelled
	})).Return(nil)

	s := svc.NewReservationService(rsvRepo, nil, nil, nil, newRedis(t), trMgr)
	err := s.ConfirmReservation(context.Background(), "rsv-001")

	assert.NotNil(t, err)
	assert.Equal(t, apperror.CategoryUnprocessable, err.Category)
}
