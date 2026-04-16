package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/model"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id uint) (*model.Subscription, error)
	GetAll() ([]model.Subscription, error)
	Update(id uint, sub *model.Subscription) error
	Delete(id uint) error
	GetTotalCost(userID *uuid.UUID, serviceName *string, startPeriod, endPeriod time.Time) (int, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *subscriptionRepo) GetByID(id uint) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.First(&sub, id).Error
	return &sub, err
}

func (r *subscriptionRepo) GetAll() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Find(&subs).Error
	return subs, err
}

func (r *subscriptionRepo) Update(id uint, sub *model.Subscription) error {
	// Updates обновляет только измененные поля
	return r.db.Model(&model.Subscription{}).Where("id = ?", id).Updates(sub).Error
}

func (r *subscriptionRepo) Delete(id uint) error {
	return r.db.Delete(&model.Subscription{}, id).Error
}

func (r *subscriptionRepo) GetTotalCost(userID *uuid.UUID, serviceName *string, startPeriod, endPeriod time.Time) (int, error) {
	var total int
	query := r.db.Model(&model.Subscription{}).
		Where("start_date >= ? AND start_date <= ?", startPeriod, endPeriod)

	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != nil && *serviceName != "" {
		query = query.Where("service_name = ?", *serviceName)
	}

	err := query.Select("COALESCE(SUM(price), 0)").Scan(&total).Error
	return total, err
}