package rest

import (
	"github.com/Yoshikrit/reservation/config"
	"github.com/Yoshikrit/reservation/internal/controller/rest/middleware"
	productCtrl "github.com/Yoshikrit/reservation/internal/controller/rest/v1/product"
	reservationCtrl "github.com/Yoshikrit/reservation/internal/controller/rest/v1/reservation"
	inventoryGateway "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory"
	inventorypb "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory/pb"
	cfgPromotionRepo "github.com/Yoshikrit/reservation/internal/repository/cfgpromotion"
	outboxRepo "github.com/Yoshikrit/reservation/internal/repository/outbox"
	reservationRepo "github.com/Yoshikrit/reservation/internal/repository/reservation"
	reservationSrv "github.com/Yoshikrit/reservation/internal/service/reservation"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	trm "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func NewRestRouter(router *fiber.App, db *gorm.DB, rdb *redis.Client, grpcConns *config.GrpcConns) {
	for _, m := range middleware.NewMiddleware() {
		router.Use(m)
	}

	newHealthRouter(router, db)
	newMetricsRouter(router)

	v1 := router.Group("/api/v1")
	newRoute(v1, db, rdb, grpcConns)
}

func newRoute(v1 fiber.Router, db *gorm.DB, rdb *redis.Client, grpcConns *config.GrpcConns) {
	trManager, err := trm.New(gormtrm.NewDefaultFactory(db))
	if err != nil {
		log.Fatal().Err(err).Msg("rest: failed to create transaction manager")
	}

	reservationRepo := reservationRepo.NewReservationRepository(db)
	cfgPromotionRepo := cfgPromotionRepo.NewCfgPromotionDateRepository(db)
	outboxRepo := outboxRepo.NewOutboxRepository(db)
	inventoryGateway := inventoryGateway.NewInventoryGateway(inventorypb.NewInventoryServiceClient(grpcConns.Inventory))
	reservationSvc := reservationSrv.NewReservationService(reservationRepo, cfgPromotionRepo, outboxRepo, inventoryGateway, rdb, trManager)

	productCtrl.NewProductController(reservationSvc).RegisterRoutes(v1.Group("/products"))
	reservationCtrl.NewReservationController(reservationSvc).RegisterRoutes(v1.Group("/reservations"))
}
