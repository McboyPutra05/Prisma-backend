package models

import "time"

type Pembayaran struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CustomerID string    `gorm:"type:varchar(50);index" json:"customer_id"`
	Bulan      string    `gorm:"type:varchar(50)" json:"bulan"`
	Tanggal    string    `gorm:"type:varchar(20)" json:"tanggal"`
	Metode     string    `gorm:"type:varchar(50)" json:"metode"`
	Nominal    float64   `json:"nominal"`
	Keterangan string    `gorm:"type:text" json:"keterangan"`
	Status     string    `gorm:"type:varchar(50);default:'Verified'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}