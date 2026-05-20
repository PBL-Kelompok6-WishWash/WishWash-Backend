package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware menginisialisasi middleware untuk autentikasi rute berbasis JSON Web Token (JWT).
// Fungsi ini memvalidasi keberadaan, format, dan integritas token sebelum meneruskan request ke Controller.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ekstraksi header Authorization dari HTTP Request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak: Header Authorization tidak ditemukan"})
			c.Abort()
			return
		}

		// 2. Validasi format Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid. Gunakan format: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "rahasia_wishwash_pbl_6" 
		}
		
		// 3. Parsing dan verifikasi signature token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi algoritma signing untuk mencegah serangan downgrade (misal: "none" algorithm)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak terduga: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		// 4. Penanganan token tidak valid atau expired
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau telah kedaluwarsa"})
			c.Abort()
			return
		}

		// 5. Ekstraksi payload (claims) dan injeksi data ke dalam context GIN
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id_user", claims["id_user"])
			c.Set("id_role", claims["id_role"])
			c.Set("username", claims["username"])
			
			// Lanjutkan eksekusi ke handler berikutnya (Controller)
			c.Next() 
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal mengekstraksi klaim dari token"})
			c.Abort()
			return
		}
	}
}

// AdminOnly adalah middleware untuk memastikan hanya pengguna dengan role Admin (RoleID = 1) yang bisa mengakses rute.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil data id_role yang sudah disimpan oleh JWTAuthMiddleware
		roleData, exists := c.Get("id_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Data role tidak ditemukan di sesi ini"})
			c.Abort()
			return
		}

		// JWT MapClaims selalu mengubah angka menjadi float64 secara default
		roleID := int(roleData.(float64))

		// Cek apakah dia Admin (1)
		if roleID != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! 🛑 Rute ini khusus untuk Admin."})
			c.Abort()
			return
		}

		// Kalau lolos, lanjut ke tujuan
		c.Next()
	}
}

// KaryawanOnly adalah middleware untuk memastikan hanya Karyawan (RoleID = 2) yang bisa mengakses.
func KaryawanOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleData, exists := c.Get("id_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Data role tidak ditemukan di sesi ini"})
			c.Abort()
			return
		}

		roleID := int(roleData.(float64))

		// Hanya Karyawan (2) yang boleh masuk
		if roleID != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! 🛑 Rute ini khusus untuk Karyawan."})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PelangganOnly adalah middleware untuk memastikan hanya Pelanggan (RoleID = 3) yang bisa mengakses.
func PelangganOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleData, exists := c.Get("id_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Data role tidak ditemukan di sesi ini"})
			c.Abort()
			return
		}

		roleID := int(roleData.(float64))

		// Cek apakah dia Pelanggan (3)
		if roleID != 3 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak! 🛑 Rute ini khusus untuk Pelanggan."})
			c.Abort()
			return
		}

		c.Next()
	}
}