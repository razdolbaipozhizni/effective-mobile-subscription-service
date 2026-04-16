package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceName string     `gorm:"column:service_name;type:varchar(255);not null" json:"service_name"`
	Price       int        `gorm:"not null" json:"price"`
	UserID      uuid.UUID  `gorm:"column:user_id;type:uuid;not null" json:"user_id"` // Возвращаем uuid.UUID
	StartDate   time.Time  `gorm:"column:start_date;not null" json:"start_date"`
	EndDate     *time.Time `gorm:"column:end_date" json:"end_date,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}