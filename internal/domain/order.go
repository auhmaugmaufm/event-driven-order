package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	TotalAmount uint      `json:"total_amount" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	User  *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
