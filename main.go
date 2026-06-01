package main

import (
	"log"
	"os"
	"time"

	"prisma-laundry-backend/config"
	"prisma-laundry-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Membaca file .env (Password dan settingan DB)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal membaca file .env! Pastikan file .env sudah ada.")
	}

	// 2. Konek ke Database MySQL
	config.ConnectDatabase()

	// 3. Inisialisasi Router Gin
	r := gin.Default()

	// 4. Setup CORS (SANGAT PENTING!)
	// Ini agar Frontend React (misal jalan di localhost:5173) diizinkan 
	// mengambil data dari Backend Golang (localhost:8080)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"}, // Sesuaikan dengan port React kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 5. Daftarkan semua Routes (URL API)
	routes.SetupRoutes(r)

	// 6. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port kalau di .env tidak ditulis
	}
	
	log.Printf("🚀 Server Prisma Laundry berjalan di http://localhost:%s", port)
	r.Run(":" + port)
}