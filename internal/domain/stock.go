package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stock struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type StockRepository interface {
	CreateWithTx(tx *gorm.DB, stock *Stock) error
	IncreaseStock(productId uuid.UUID, quantity int) error
	DecreaseStock(productId uuid.UUID, quantity int) error
	GetProductStock(productId uuid.UUID) (*Stock, error)
	GetStocks() ([]Stock, error)
}
