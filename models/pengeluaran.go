package models

import "time"

type Pengeluaran struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Bulan     string    `gorm:"type:varchar(50)" json:"bulan"`
	Tanggal   string    `gorm:"type:varchar(20)" json:"tanggal"`
	Barang    string    `gorm:"type:text" json:"barang"` 
	Total     string    `gorm:"type:varchar(50)" json:"total"`
	
	CreatedAt time.Time `json:"created_at"`
}