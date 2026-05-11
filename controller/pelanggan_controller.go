package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
)

// DTO untuk input dari Next.js
type PelangganInput struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password"` // Bisa kosong saat update
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	NoTelp      string `json:"no_telp"`
	FotoPelanggan string `json:"foto_pelanggan"`
}

type PelangganController interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type pelangganController struct {
	pelangganRepo repository.PelangganRepository
	userRepo      repository.UserRepository // Butuh ini untuk insert/hapus akun User-nya juga
}

func NewPelangganController(pRepo repository.PelangganRepository, uRepo repository.UserRepository) PelangganController {
	return &pelangganController{pRepo, uRepo}
}

// Helper untuk konversi string (dari URL) ke uint
func parseID(idStr string) uint {
	id, _ := strconv.ParseUint(idStr, 10, 32)
	return uint(id)
}

func (ctrl *pelangganController) GetAll(c *gin.Context) {
	pelanggans, err := ctrl.pelangganRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pelanggan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pelanggans})
}

func (ctrl *pelangganController) GetByID(c *gin.Context) {
	id := parseID(c.Param("id"))

	// Cari pelanggan berdasarkan ID (Sudah otomatis Preload User di Repo)
	pelanggan, err := ctrl.pelangganRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pelanggan})
}

func (ctrl *pelangganController) Create(c *gin.Context) {
	var input PelangganInput
	
	// 1. Validasi Input (Tangkap Error Spesifik)
	if err := c.ShouldBindJSON(&input); err != nil {
		var errorMsgs []string

		// Cek apakah errornya dari tag 'binding' (validator)
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				switch fieldErr.Field() {
				case "Username":
					errorMsgs = append(errorMsgs, "Username tidak boleh kosong")
				case "Email":
					if fieldErr.Tag() == "required" {
						errorMsgs = append(errorMsgs, "Email tidak boleh kosong")
					} else if fieldErr.Tag() == "email" {
						errorMsgs = append(errorMsgs, "Format email tidak valid (harus ada @ dsb)")
					}
				case "NamaLengkap":
					errorMsgs = append(errorMsgs, "Nama Lengkap tidak boleh kosong")
				}
			}
			// Gabungkan semua pesan error jadi satu kalimat
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMsgs, ", ")})
			return
		}

		// Kalau errornya karena format JSON rusak
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data JSON tidak valid"})
		return
	}

	// Validasi Manual Password (Karena saat Update boleh kosong, tapi saat Create WAJIB ada)
	if input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak boleh kosong"})
		return
	}

	// 2. Buat Akun User Dulu
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   3, // 3 = Role Pelanggan
	}

	// Simpan User
	if err := ctrl.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun user"})
		return
	}

	// 3. Buat Profil Pelanggan
	pelanggan := model.Pelanggan{
		UserID:        user.IDUser,
		NamaLengkap:   input.NamaLengkap,
		NoTelp:        input.NoTelp,
		FotoPelanggan: input.FotoPelanggan, // 💡 Pastikan foto ikut disimpan ke DB
	}

	// Simpan Profil Pelanggan
	if err := ctrl.pelangganRepo.CreatePelanggan(&pelanggan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan detail pelanggan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pelanggan berhasil ditambahkan!", "data": pelanggan})
}

func (ctrl *pelangganController) Update(c *gin.Context) {
	id := parseID(c.Param("id"))
	var input PelangganInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errorMsgs []string
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				switch fieldErr.Field() {
				case "Username":
					errorMsgs = append(errorMsgs, "Username tidak boleh kosong")
				case "Email":
					if fieldErr.Tag() == "required" {
						errorMsgs = append(errorMsgs, "Email tidak boleh kosong")
					} else if fieldErr.Tag() == "email" {
						errorMsgs = append(errorMsgs, "Format email tidak valid")
					}
				case "NamaLengkap":
					errorMsgs = append(errorMsgs, "Nama Lengkap tidak boleh kosong")
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMsgs, ", ")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	pelanggan, err := ctrl.pelangganRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	pelanggan.User.Username = input.Username
	pelanggan.User.Email = input.Email

	// Cek jika password diisi, baru kita hash & update. Kalau kosong, biarin password lama.
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		pelanggan.User.Password = string(hashedPassword)
	}

	if err := ctrl.userRepo.UpdateUser(&pelanggan.User); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate akun user"})
		return
	}

	pelanggan.NamaLengkap = input.NamaLengkap
	pelanggan.NoTelp = input.NoTelp
	pelanggan.FotoPelanggan = input.FotoPelanggan // 💡 Simpan fotonya juga

	if err := ctrl.pelangganRepo.Update(pelanggan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data pelanggan berhasil diperbarui!"})
}

func (ctrl *pelangganController) Delete(c *gin.Context) {
	id := parseID(c.Param("id"))

	// 1. Cari data pelanggan untuk mendapatkan UserID
	pelanggan, err := ctrl.pelangganRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	// 2. Hapus Profil Pelanggannya dulu (Child)
	if err := ctrl.pelangganRepo.Delete(pelanggan.IDPelanggan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus profil pelanggan"})
		return
	}

	// 3. Hapus Akun User-nya (Parent) agar tidak jadi sampah
	// Meminjam koneksi DB dari userRepo untuk menghapus berdasarkan UserID
	if err := ctrl.userRepo.DeleteUser(pelanggan.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus akun user pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Pelanggan berhasil dihapus"})
}