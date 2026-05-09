package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 1. DTO (Data Transfer Object) disesuaikan dengan model.User yang asli
type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
    NoTelp      string `json:"no_telp" binding:"required"`
	RoleID   uint   `json:"id_role" binding:"required"` // Pakai ID (angka), bukan string teks
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authController struct {
	userRepo repository.UserRepository
	pelangganRepo repository.PelangganRepository
	karyawanRepo  repository.KaryawanRepository
}

func NewAuthController(userRepo repository.UserRepository, 
					   pelangganRepo repository.PelangganRepository,
					   karyawanRepo repository.KaryawanRepository,) AuthController {
	return &authController{userRepo, pelangganRepo, karyawanRepo}
}

// 4. Logika Register
func (ctrl *authController) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap atau format salah"})
		return
	}

	if input.RoleID != 1 && input.RoleID != 2 && input.RoleID != 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid! Masukkan 1 (Karyawan), 2 (Pelanggan), atau 3 (Admin)"})
		return
	}

	// 1. Cek Username (Pencegatan yang kita buat sebelumnya)
	_, err := ctrl.userRepo.FindByUsername(input.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar!"})
		return
	}

	// 2. Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// 3. Siapkan Model User
	user := model.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
		RoleID:   input.RoleID,
	}

	// 4. Simpan ke tabel 'user'
	if err := ctrl.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah digunakan!"})
		return
	}

	// 5. SIMPAN KE TABEL PELANGGAN ATAU KARYAWAN BERDASARKAN ROLE
	if input.RoleID == 3 {
		pelanggan := model.Pelanggan{
			UserID:      user.IDUser,
			NamaLengkap: input.NamaLengkap,
			NoTelp:      input.NoTelp,
		}
		if err := ctrl.pelangganRepo.CreatePelanggan(&pelanggan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil pelanggan"})
			return
		}
	} else if input.RoleID == 2 {
		karyawan := model.Karyawan{
			UserID:             user.IDUser,
			NamaKaryawan:       input.NamaLengkap, // Kita pakai input.NamaLengkap untuk mengisi NamaKaryawan
			NoTelp:             input.NoTelp,
			StatusKetersediaan: "Tersedia",        // Beri nilai default
		}
		if err := ctrl.karyawanRepo.CreateKaryawan(&karyawan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil karyawan"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Registrasi akun berhasil!",
		"username": user.Username,
	})
}

// 5. Logika Login
func (ctrl *authController) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userRepo.FindByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "rahasia_wishwash_pbl_6" 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_user":  user.IDUser,
		"username": user.Username,
		"id_role":  user.RoleID, // Menggunakan user.RoleID sesuai database
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
		"id_role": user.RoleID,
	})
}