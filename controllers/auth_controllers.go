package controllers

import (
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"
	"prisma-laundry-backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input models.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah!"})
		return
	}

	var user models.User

	if err := config.DB.Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No. WhatsApp tidak terdaftar!"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat sesi login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Berhasil!",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

type InputUser struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"` 
}

func CreateUser(c *gin.Context) {
	var input InputUser

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah. Pastikan name, phone, password, dan role terisi."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No. WA mungkin sudah terdaftar atau gagal menyimpan ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Akun " + input.Name + " dengan role " + input.Role + " berhasil dibuat!",
	})
}