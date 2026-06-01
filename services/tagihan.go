package services

import (
    "prisma-laundry-backend/config"
    "prisma-laundry-backend/models"
)

type TagihanBulananRes struct {
	CustomerName string `json:"customer_name"`
	Jan          int    `json:"jan"`
	Feb          int    `json:"feb"`
	Mar          int    `json:"mar"`
	Apr          int    `json:"apr"`
	May          int    `json:"may"`
	Jun          int    `json:"jun"`
	Jul          int    `json:"jul"`
	Aug          int    `json:"aug"`
	Sep          int    `json:"sep"`
	Oct          int    `json:"oct"`
	Nov          int    `json:"nov"`
	Dec          int    `json:"dec"`
}

func CreateTagihan(tagihan *models.Tagihan) error {
	total := 0

	for _, item := range tagihan.Details {
		total += item.Harga
	}

	tagihan.Total = total
	tagihan.Status = "pending"

	return config.DB.Create(tagihan).Error
}

func GetTagihanByUser(userID uint) ([]models.Tagihan, error) {
	var tagihan []models.Tagihan

	err := config.DB.
		Preload("Details").
		Where("user_id = ?", userID).
		Find(&tagihan).Error

	return tagihan, err
}

func UpdateStatusTagihan(id uint, status string) error {
	return config.DB.Model(&models.Tagihan{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Tambahkan fungsi ini di paling bawah file services/tagihan.go
func GetRekapTagihanSemuaCustomer() ([]TagihanBulananRes, error) {
	var tagihans []models.Tagihan
	
	// Preload User untuk mendapatkan nama Customer
	err := config.DB.Preload("User").Find(&tagihans).Error
	if err != nil {
		return nil, err
	}

	// Gunakan map untuk mengelompokkan data berdasarkan UserID
	rekapMap := make(map[uint]*TagihanBulananRes)

	for _, t := range tagihans {
		if _, exists := rekapMap[t.UserID]; !exists {
			rekapMap[t.UserID] = &TagihanBulananRes{
				CustomerName: t.User.Name, // Asumsi di models.User ada kolom Name
			}
		}

		rekap := rekapMap[t.UserID]
		
		// Deteksi bulan dari CreatedAt dan tambahkan Total-nya
		switch t.CreatedAt.Month() {
		case 1:  rekap.Jan += t.Total
		case 2:  rekap.Feb += t.Total
		case 3:  rekap.Mar += t.Total
		case 4:  rekap.Apr += t.Total
		case 5:  rekap.May += t.Total
		case 6:  rekap.Jun += t.Total
		case 7:  rekap.Jul += t.Total
		case 8:  rekap.Aug += t.Total
		case 9:  rekap.Sep += t.Total
		case 10: rekap.Oct += t.Total
		case 11: rekap.Nov += t.Total
		case 12: rekap.Dec += t.Total
		}
	}

	// Ubah map menjadi slice/array untuk dikirim sebagai JSON
	var result []TagihanBulananRes
	for _, r := range rekapMap {
		result = append(result, *r)
	}

	return result, nil
}