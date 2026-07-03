package reservation

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

func (r *reservationRepository) Create(ctx context.Context, reservation *entity.Reservation) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperror.NewError(40900000, err, "reservation", reservation.ReservationID)
		}
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *reservationRepository) CreateBulk(ctx context.Context, reservations []entity.Reservation) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(&reservations).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *entity.Reservation) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Save(reservation).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *reservationRepository) Patch(ctx context.Context, reservation *entity.Reservation) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Model(&entity.Reservation{}).Where("reservation_id = ?", reservation.ReservationID).Updates(reservation).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *reservationRepository) Find(ctx context.Context, filter *entity.Reservation) (*entity.Reservation, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	primaryFilter, appErr := r.buildPrimaryFilter(filter)
	if appErr != nil {
		return nil, appErr
	}

	var reservation entity.Reservation
	if err := db.WithContext(ctx).Where(primaryFilter).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(40400000, err, "reservation", primaryFilter["reservation_id"])
		}
		return nil, apperror.NewError(50000000, err)
	}
	return &reservation, nil
}

func (r *reservationRepository) Filter(ctx context.Context, filter *entity.Reservation, limit, offset int, isAsc bool) ([]entity.Reservation, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	query := db.WithContext(ctx).Model(&entity.Reservation{})
	if filter != nil {
		query = query.Where(filter)
	}

	query, appErr := r.applyListOptions(query, &entity.Reservation{}, limit, offset, isAsc)
	if appErr != nil {
		return nil, appErr
	}

	var reservations []entity.Reservation
	if err := query.Find(&reservations).Error; err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return reservations, nil
}

func (r *reservationRepository) FindAll(ctx context.Context) ([]entity.Reservation, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	var reservations []entity.Reservation
	if err := db.WithContext(ctx).Find(&reservations).Error; err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return reservations, nil
}

func (r *reservationRepository) Delete(ctx context.Context, filter *entity.Reservation) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Where(filter).Delete(&entity.Reservation{}).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *reservationRepository) Count(ctx context.Context, filter *entity.Reservation) (int64, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	var count int64
	if err := db.WithContext(ctx).Model(&entity.Reservation{}).Where(filter).Count(&count).Error; err != nil {
		return 0, apperror.NewError(50000000, err)
	}
	return count, nil
}

func (r *reservationRepository) applyListOptions(query *gorm.DB, model any, limit, offset int, isAsc bool) (*gorm.DB, *apperror.AppError) {
	if limit < 0 {
		return nil, apperror.NewError(40000000, errors.New("limit must be greater than or equal to 0"))
	}
	if offset < 0 {
		return nil, apperror.NewError(40000000, errors.New("offset must be greater than or equal to 0"))
	}

	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(model); err != nil {
		return nil, apperror.NewError(50000000, err)
	}

	if len(stmt.Schema.PrimaryFields) > 0 {
		direction := "DESC"
		if isAsc {
			direction = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", stmt.Schema.PrimaryFields[0].DBName, direction))
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	return query, nil
}

func (r *reservationRepository) buildPrimaryFilter(filter *entity.Reservation) (map[string]any, *apperror.AppError) {
	if filter == nil {
		return nil, apperror.NewError(40000000, errors.New("find filter is required"))
	}

	stmt := &gorm.Statement{DB: r.db}
	if err := stmt.Parse(&entity.Reservation{}); err != nil {
		return nil, apperror.NewError(50000000, err)
	}

	value := reflect.ValueOf(filter)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, apperror.NewError(40000000, errors.New("find filter is required"))
		}
		value = value.Elem()
	}

	primaryFilter := make(map[string]any, len(stmt.Schema.PrimaryFields))
	for _, primaryField := range stmt.Schema.PrimaryFields {
		fieldValue := value.FieldByName(primaryField.Name)
		if !fieldValue.IsValid() || fieldValue.IsZero() {
			return nil, apperror.NewError(40000000, fmt.Errorf("missing primary field: %s", primaryField.DBName))
		}
		primaryFilter[primaryField.DBName] = fieldValue.Interface()
	}

	return primaryFilter, nil
}
