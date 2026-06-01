package controllers

import (
	"fmt"
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"
	"prisma-laundry-backend/utils"

	"github.com/gin-gonic/gin"
)

// 1. Fungsi Tambah Pengeluaran (Ditembak saat klik "Selesai" di Modal Bon)
func CreatePengeluaran(c *gin.Context) {
	var input models.Pengeluaran

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengeluaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data bon berhasil disimpan!", "data": input})
}

// 2. Fungsi Ambil Semua Pengeluaran (Bisa difilter berdasarkan bulan)
func GetPengeluaran(c *gin.Context) {
	var pengeluaran []models.Pengeluaran
	
	bulan := c.Query("bulan")

	query := config.DB
	if bulan != "" {
		query = query.Where("bulan = ?", bulan)
	}

	if err := query.Order("created_at desc").Find(&pengeluaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pengeluaran})
}

func ExportPengeluaran(c *gin.Context) {
	bulan := c.Query("bulan")

	if bulan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter bulan wajib diisi"})
		return
	}

	var pengeluarans []models.Pengeluaran
	
	// 1. Ambil Data dari Database
	if err := config.DB.Where("bulan = ?", bulan).Order("tanggal asc").Find(&pengeluarans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengeluaran dari database"})
		return
	}

	// 2. Lempar data ke file utils untuk dirakit jadi Excel
	fileExcel, err := utils.GeneratePengeluaranExcel(pengeluarans)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merakit data Excel"})
		return
	}

	// 3. Kirim file yang sudah jadi ke Browser
	filename := fmt.Sprintf("Data_Pengeluaran_%s.xlsx", bulan)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if err := fileExcel.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim file Excel"})
		return
	}
}