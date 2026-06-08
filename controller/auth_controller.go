package controller

import (
	"fmt"
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
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Email       string `json:"email" binding:"required,email"`
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	NoTelp      string `json:"no_telp" binding:"required"`
	RoleID      uint   `json:"id_role" binding:"required"` // Pakai ID (angka), bukan string teks
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
	userRepo       repository.UserRepository
	pelangganRepo  repository.PelangganRepository
	karyawanRepo   repository.KaryawanRepository
	adminRepo      repository.AdminRepository
	notifikasiRepo repository.NotifikasiRepository
}

func NewAuthController(userRepo repository.UserRepository,
						pelangganRepo repository.PelangganRepository,
						karyawanRepo repository.KaryawanRepository,
						adminRepo repository.AdminRepository,
						notifikasiRepo repository.NotifikasiRepository) AuthController {
						return &authController{userRepo, pelangganRepo, karyawanRepo, adminRepo, notifikasiRepo}
}

// 4. Logika Register
func (ctrl *authController) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap atau format salah"})
		return
	}

	if input.RoleID != 1 && input.RoleID != 2 && input.RoleID != 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid! Masukkan 1 (Admin), 2 (Karyawan), atau 3 (Pelanggan)"})
		return
	}

	if input.RoleID != 2 && input.RoleID != 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak valid atau tidak diizinkan! Admin hanya bisa dibuat oleh sistem. 🛑"})
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
	switch input.RoleID {
	case 3:
		pelanggan := model.Pelanggan{
			UserID:      user.IDUser,
			NamaLengkap: input.NamaLengkap,
			NoTelp:      input.NoTelp,
		}
		if err := ctrl.pelangganRepo.CreatePelanggan(&pelanggan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil pelanggan"})
			return
		}
		
		// Trigger notification for admins
		go ctrl.notifikasiRepo.CreateNotificationForAdmins("Pelanggan Baru 🎉", fmt.Sprintf("Pelanggan baru bernama %s (@%s) telah terdaftar.", pelanggan.NamaLengkap, user.Username))
	case 2:
		karyawan := model.Karyawan{
			UserID:             user.IDUser,
			NamaKaryawan:       input.NamaLengkap, // Kita pakai input.NamaLengkap untuk mengisi NamaKaryawan
			NoTelp:             input.NoTelp,
			StatusKetersediaan: "Tersedia", // Beri nilai default
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format permintaan tidak valid."})
		return
	}

	user, err := ctrl.userRepo.FindByUsername(input.Username)
	if err != nil {
		// Pesan error spesifik untuk Username
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username yang Anda masukkan tidak terdaftar di sistem."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// Pesan error spesifik untuk Password
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kata sandi yang Anda masukkan tidak sesuai."})
		return
	}

	var displayName string

	switch user.RoleID {
	case 1: // ADMIN
		// Admin tetap pakai userRepo jika kamu belum buat adminRepo
		admin, err := ctrl.adminRepo.FindByUserID(user.IDUser)
		if err == nil {
			displayName = admin.NamaAdmin
		}
	case 2: // KARYAWAN
		// 💡 Panggil dari karyawanRepo
		karyawan, err := ctrl.karyawanRepo.FindByUserID(user.IDUser)
		if err == nil {
			displayName = karyawan.NamaKaryawan
		}
	case 3: // PELANGGAN
		// 💡 Panggil dari pelangganRepo
		pelanggan, err := ctrl.pelangganRepo.FindByUserID(user.IDUser)
		if err == nil {
			displayName = pelanggan.NamaLengkap
		}
	default:
		displayName = "User"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan sistem saat membuat sesi autentikasi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Autentikasi berhasil.",
		"token":        tokenString,
		"id_role":      user.RoleID,
		"id_user":      user.IDUser,
		"display_name": displayName,
	})
}
