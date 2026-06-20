package models

import "time"

type StatusBarang struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CustomerID   uint      `json:"customer_id"` 
	CustomerName string    `gorm:"type:varchar(150)" json:"customer_name"`
	Tanggal      string    `gorm:"type:varchar(20)" json:"tanggal"`
	Code         string    `gorm:"type:varchar(50)" json:"code"`           // cth: "PO 77"
	JenisBarang  string    `gorm:"type:varchar(150)" json:"jenis_barang"`  // cth: "Celana Pendek Dws"
	Qty          string    `gorm:"type:varchar(100)" json:"qty"`           // cth: "70 Lusin"
	JenisCucian  string    `gorm:"type:text" json:"jenis_cucian"`          // cth: "Bio scrub w"
	Keterangan   string    `gorm:"type:text" json:"keterangan"` // Akan menyimpan: [{"qty":"15", "type":"Bio"}]
    Status       string    `gorm:"type:enum('Belum diproses','Diproses','Selesai');default:'Belum diproses'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}