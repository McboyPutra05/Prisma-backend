package models

import "time"

type Obat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Bulan     string    `gorm:"type:varchar(50)" json:"bulan"`     // cth: "Februari 2026"
	Tanggal   string    `gorm:"type:varchar(20)" json:"tanggal"`   // cth: "2026-02-15"
	Name      string    `gorm:"type:varchar(255)" json:"name"`     // cth: "Softener Flake"
	Qty       string    `gorm:"type:varchar(50)" json:"qty"`       // cth: "50 kg"
	CreatedAt time.Time `json:"created_at"`
}