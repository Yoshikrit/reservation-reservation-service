package reservation

import (
	"context"

	inventoryGateway "reservation/internal/gateway/grpc/inventory"
	"reservation/internal/pkg/apperror"
	cfgPromotionRepo "reservation/internal/repository/cfgpromotion"
	outboxRepo "reservation/internal/repository/outbox"
	reservationRepo "reservation/internal/repository/reservation"

	"github.com/redis/go-redis/v9"
)

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type ReservationService interface {
	CreateProduct(ctx context.Context, request *CreateProductRequest) *apperror.AppError
	GetProduct(ctx context.Context, productID string) (*GetProductResponse, *apperror.AppError)
	CreateReservation(ctx context.Context, request *CreateReservationRequest) *apperror.AppError
	GetReservations(ctx context.Context, body *ListReservationBody) (*ListReservationResponse, *apperror.AppError)
	CancelReservation(ctx context.Context, reservationID string) *apperror.AppError
	ConfirmReservation(ctx context.Context, reservationID string) *apperror.AppError
}

type reservationService struct {
	inventoryGateway inventoryGateway.InventoryGateway
	reservationRepo  reservationRepo.ReservationRepository
	cfgPromotionRepo cfgPromotionRepo.CfgPromotionDateRepository
	outboxRepo       outboxRepo.OutboxRepository
	redis            *redis.Client
	trManager        TrManager
}

func NewReservationService(
	reservationRepo reservationRepo.ReservationRepository,
	cfgPromotionRepo cfgPromotionRepo.CfgPromotionDateRepository,
	outboxRepo outboxRepo.OutboxRepository,
	inventoryGateway inventoryGateway.InventoryGateway,
	rdb *redis.Client,
	trManager TrManager,
) ReservationService {
	return &reservationService{
		reservationRepo:  reservationRepo,
		cfgPromotionRepo: cfgPromotionRepo,
		outboxRepo:       outboxRepo,
		inventoryGateway: inventoryGateway,
		redis:            rdb,
		trManager:        trManager,
	}
}
