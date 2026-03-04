package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductName  string    `json:"product_name" gorm:"not null"`
	ProductPrice uint      `json:"product_price" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
