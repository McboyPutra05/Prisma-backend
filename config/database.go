package config

import (
	"fmt"
	"log"
	"os"

	"prisma-laundry-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal terhubung ke database MySQL: ", err)
	}

	err = database.AutoMigrate(
		&models.User{},
		&models.Obat{},
		&models.Pengeluaran{},
		&models.StatusBarang{},
		&models.Tagihan{},
		&models.TagihanRinci{},
		&models.Customer{},
		&models.Pembayaran{},
	)

	if err != nil {
		log.Fatal("❌ Gagal melakukan AutoMigrate: ", err)
	}

	DB = database
	fmt.Println("✅ Berhasil terhubung ke database MySQL dan AutoMigrate sukses!")
}