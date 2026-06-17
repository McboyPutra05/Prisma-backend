package controllers

import (
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	customerID := c.Param("id")
	var notifications []models.Notification

	// Ambil notifikasi berdasarkan customer_id (atau "PETUGAS"), urutkan dari yang terbaru
	if err := config.DB.Where("customer_id = ?", customerID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}

func MarkNotificationRead(c *gin.Context) {
	id := c.Param("id")
	var notification models.Notification

	if err := config.DB.First(&notification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan"})
		return
	}

	if err := config.DB.Model(&notification).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi telah dibaca"})
}
