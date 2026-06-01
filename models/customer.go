package models

import "time"

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(150);not null" json:"name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Address   string    `gorm:"type:text" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}