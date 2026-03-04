package domain

import "github.com/google/uuid"

type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     uint      `json:"price" gorm:"not null"`

	Product Product `json:"product" gorm:"foreignKey:ProductID"`
	Order   Order   `json:"order" gorm:"foreignKey:OrderID"`
}
