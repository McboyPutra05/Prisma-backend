package models

import "gorm.io/gorm"

type Tagihan struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID"`

	Total  int    `json:"total"`
	Status string `json:"status"` // pending, lunas

	Details []TagihanDetail `gorm:"foreignKey:TagihanID"`
}

type TagihanDetail struct {
	gorm.Model
	TagihanID uint   `json:"tagihan_id"`
	No        string `json:"no"` 
	Tgl       string `json:"tgl"` // Bisa berisi tanggal (03/02/2026) atau PO (J 319)
	Qty       string `json:"qty"`
	NamaItem  string `json:"nama_item"` // Item (c panjang dws)
	Trt       string `json:"trt"`       // Treatment (bio scrb w)
	Harga     int    `json:"harga"`
	TotalRow  int    `json:"total_row"` // Total per baris (Qty * Harga)

	Highlight string `json:"highlight"` 
}