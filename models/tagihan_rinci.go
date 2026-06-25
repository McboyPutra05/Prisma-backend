package models

import "time"

type TagihanRinci struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CustomerID  string    `gorm:"type:varchar(50);index" json:"customer_id"` 
	Bulan       string    `gorm:"type:varchar(50)" json:"bulan"`
	Tanggal     string    `gorm:"type:varchar(20)" json:"tanggal"`
	Merk        string    `gorm:"type:varchar(100)" json:"merk"`
	KodePo      string    `gorm:"type:varchar(100)" json:"kode_po"`
	Qty         string    `gorm:"type:varchar(20)" json:"qty"`
	Item        string    `gorm:"type:varchar(255)" json:"item"`
	JenisCucian string    `gorm:"type:varchar(255)" json:"jenis_cucian"`
	NoCelana    string    `gorm:"type:varchar(50)" json:"no_celana"`
	Harga       float64   `json:"harga"`
	Total       float64   `json:"total"`
	Status      string    `gorm:"type:varchar(50);default:'Belum Lunas'" json:"status"` 
	
	CreatedAt   time.Time `json:"created_at"`
}