package controllers

import (
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCustomersDropdown(c *gin.Context) {
	var customers []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	config.DB.Table("users").Select("id, name").Where("role = 'customer'").Scan(&customers)
	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func CreatePembayaran(c *gin.Context) {
	var input models.Pembayaran
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bulan := "Januari"
	if len(input.Tanggal) >= 7 {
		monthStr := input.Tanggal[5:7]
		monthInt, _ := strconv.Atoi(monthStr)
		daftarBulan := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
		if monthInt >= 1 && monthInt <= 12 {
			bulan = daftarBulan[monthInt]
		}
	}
	input.Bulan = bulan
	input.Status = "Verified"

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil disimpan!"})
}

func GetPembayaran(c *gin.Context) {
	bulan := c.Query("bulan")
	var pembayaranList []models.Pembayaran

	query := config.DB.Model(&models.Pembayaran{})
	if bulan != "" && bulan != "Pilih Bulan" {
		query = query.Where("bulan = ?", bulan)
	}
	query.Order("tanggal desc").Find(&pembayaranList)

	var totalNominal float64 = 0
	var formattedDetails []map[string]interface{}

	for _, p := range pembayaranList {
		var customer struct{ Name string }
		config.DB.Table("users").Select("name").Where("id = ?", p.CustomerID).Scan(&customer)

		formattedDetails = append(formattedDetails, map[string]interface{}{
			"id":         p.ID,
			"tgl":        p.Tanggal,
			"customer":   customer.Name,
			"keterangan": p.Keterangan,
			"metode":     p.Metode,
			"nominal":    p.Nominal,
			"status":     p.Status,
		})
		totalNominal += p.Nominal
	}

	c.JSON(http.StatusOK, gin.H{
		"data": formattedDetails,
		"summary": gin.H{
			"total_diterima":  totalNominal,
			"total_transaksi": len(pembayaranList),
		},
	})
}

func GetPembayaranCustomer(c *gin.Context) {
	customerId := c.Param("id")

	var pembayaranList []models.Pembayaran
	config.DB.Where("customer_id = ?", customerId).Order("tanggal desc").Find(&pembayaranList)

	var formattedDetails []map[string]interface{}

	for _, p := range pembayaranList {
		tahun := ""
		if len(p.Tanggal) >= 4 {
			tahun = p.Tanggal[:4]
		}

		formattedDetails = append(formattedDetails, map[string]interface{}{
			"id":     p.ID,
			"date":   p.Tanggal,
			"amount": p.Nominal,
			"desc":   p.Keterangan,
			"status": p.Status,
			"month":  p.Bulan,
			"year":   tahun,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": formattedDetails,
	})
}