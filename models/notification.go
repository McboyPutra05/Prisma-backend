package models

import "time"

type Notification struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CustomerID string    `gorm:"type:varchar(50);index" json:"customer_id"`
	Title      string    `gorm:"type:varchar(100)" json:"title"`
	Message    string    `gorm:"type:text" json:"message"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}