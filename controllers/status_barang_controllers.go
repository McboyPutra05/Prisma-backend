package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"

	"github.com/gin-gonic/gin"
)

func CreateStatusBarang(c *gin.Context) {
	var input models.StatusBarang

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data status barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data status barang berhasil ditambahkan!", "data": input})
}

func GetStatusBarang(c *gin.Context) {
	var statusBarangs []models.StatusBarang

	status := c.Query("status")

	query := config.DB
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at desc").Find(&statusBarangs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data status barang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": statusBarangs})
}

type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func UpdateStatusBarang(c *gin.Context) {
	id := c.Param("id")
	var input UpdateStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak boleh kosong"})
		return
	}
	var barang models.StatusBarang
	if err := config.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data status barang tidak ditemukan"})
		return
	}

	if err := config.DB.Model(&barang).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update status barang"})
		return
	}

	if input.Status == "Selesai" {
		notif := models.Notification{
			CustomerID: fmt.Sprintf("%v", barang.CustomerID),
			Title:      "Cucian Selesai! 🎉",
			Message:    "Hore! Cucian kamu dengan PO " + barang.Code + " sudah selesai diproses dan siap diambil/dikirim.",
			IsRead:     false,
		}
		config.DB.Create(&notif)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status barang berhasil diupdate menjadi: " + input.Status})
}

type UpdateKeteranganInput struct {
	Keterangan interface{} `json:"keterangan" binding:"required"`
}

func UpdateKeteranganBarang(c *gin.Context) {
	id := c.Param("id")
	var input UpdateKeteranganInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	keteranganBytes, _ := json.Marshal(input.Keterangan)
	keteranganString := string(keteranganBytes)

	result := config.DB.Model(&models.StatusBarang{}).Where("id = ?", id).Update("keterangan", keteranganString)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan keterangan"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data status barang tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Keterangan berhasil disimpan!"})
}

func DeleteStatusBarang(c *gin.Context) {
    id := c.Param("id")
    if err := config.DB.Delete(&models.StatusBarang{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Data barang berhasil dihapus!"})
}