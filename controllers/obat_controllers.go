package controllers

import (
	"fmt"
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"
	"prisma-laundry-backend/utils"

	"github.com/gin-gonic/gin"
)

type InputObatRequest struct {
	Bulan   string `json:"bulan"`
	Tanggal string `json:"tanggal"`
	Items   []struct {
		Name string `json:"name"`
		Qty  string `json:"qty"`
	} `json:"items"`
}

func CreateObat(c *gin.Context) {
	var input InputObatRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data obat salah!"})
		return
	}

	var savedObats []models.Obat
	for _, item := range input.Items {
		obatBaru := models.Obat{
			Bulan:   input.Bulan,
			Tanggal: input.Tanggal,
			Name:    item.Name,
			Qty:     item.Qty,
		}

		if err := config.DB.Create(&obatBaru).Error; err == nil {
			savedObats = append(savedObats, obatBaru)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data kebutuhan obat berhasil disimpan!",
		"data":    savedObats,
	})
}

// 2. Ambil Semua Data Obat (Bisa di-filter berdasarkan Bulan)
func GetObat(c *gin.Context) {
	var obats []models.Obat
	bulan := c.Query("bulan")

	query := config.DB
	if bulan != "" {
		query = query.Where("bulan = ?", bulan)
	}

	if err := query.Order("created_at desc").Find(&obats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": obats})
}

// 3. Hapus Data Obat (Berdasarkan ID)
func DeleteObat(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Obat{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data obat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data obat berhasil dihapus"})
}

func ExportObatBulanan(c *gin.Context) {
	bulan := c.Query("bulan")

	if bulan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter bulan wajib diisi"})
		return
	}

	var obats []models.Obat
	// 1. Ambil Data dari Database
	if err := config.DB.Where("bulan = ?", bulan).Order("tanggal asc").Find(&obats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data obat dari database"})
		return
	}

	// 2. Lempar data ke file utils untuk dirakit jadi Excel
	fileExcel, err := utils.GenerateObatExcel(obats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merakit data Excel"})
		return
	}

	// 3. Kirim file yang sudah jadi ke Browser
	filename := fmt.Sprintf("Data_Obat_%s.xlsx", bulan)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if err := fileExcel.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim file Excel"})
		return
	}
}