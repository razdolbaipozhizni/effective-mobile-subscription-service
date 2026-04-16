package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/model"
	"github.com/razdolbaipozhizni/effective-mobile-subscription-service/internal/repository"
)

type SubscriptionService interface {
	Create(sub *model.Subscription) error
	GetByID(id uint) (*model.Subscription, error)
	GetAll() ([]model.Subscription, error)
	Update(id uint, sub *model.Subscription) error // Добавили id
	Delete(id uint) error
	GetTotalCost(userID *uuid.UUID, serviceName *string, startPeriod, endPeriod time.Time) (int, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(sub *model.Subscription) error {
	return s.repo.Create(sub)
}

func (s *subscriptionService) GetByID(id uint) (*model.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *subscriptionService) GetAll() ([]model.Subscription, error) {
	return s.repo.GetAll()
}

func (s *subscriptionService) Update(id uint, sub *model.Subscription) error {
	return s.repo.Update(id, sub) // Передаем id в репозиторий
}

func (s *subscriptionService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *subscriptionService) GetTotalCost(userID *uuid.UUID, serviceName *string, startPeriod, endPeriod time.Time) (int, error) {
	return s.repo.GetTotalCost(userID, serviceName, startPeriod, endPeriod)
}