package controllers

import (
	"net/http"
	"prisma-laundry-backend/config"
	"prisma-laundry-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 1. Mengambil daftar pelanggan (Customer)
func GetCustomers(c *gin.Context) {
	var customers []models.User

	if err := config.DB.Where("role = ?", "customer").Order("created_at desc").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

// 2. Mengambil daftar pegawai (Petugas)
func GetPetugas(c *gin.Context) {
	var petugas []models.User

	if err := config.DB.Where("role = ?", "petugas").Order("created_at desc").Find(&petugas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data petugas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": petugas})
}

// 3. Mengambil daftar manajemen (Admin)
func GetAdmin(c *gin.Context) {
	var admin []models.User

	if err := config.DB.Where("role = ?", "admin").Order("created_at desc").Find(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admin})
}

type CreateCustomerInput struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

func CreateCustomer(c *gin.Context) {
	var input CreateCustomerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid. Pastikan semua kolom terisi."})
		return
	}

	// Enkripsi password sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	newCustomer := models.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		Role:     "customer",
	}

	if err := config.DB.Create(&newCustomer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data customer ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil ditambahkan!", "data": newCustomer})
}