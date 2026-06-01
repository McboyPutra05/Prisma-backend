package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Rahasia untuk kunci JWT (Nanti taruh di file .env ya: JWT_SECRET=rahasia_prisma_123)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Fungsi untuk membuat Token JWT saat berhasil Login
func GenerateToken(userID uint, role string) (string, error) {
	// Isi dari tiket (Payload)
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Tiket hangus dalam 24 Jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}