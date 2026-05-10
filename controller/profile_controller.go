package controller

import (
	"net/http"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// DTO yang umum untuk semua Role
type UpdateProfileInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Nama     string `json:"nama"` // Menampung NamaAdmin / NamaKaryawan / NamaLengkap
}

type ProfileController interface {
	UpdateProfile(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type profileController struct {
	userRepo      repository.UserRepository
	adminRepo     repository.AdminRepository
	karyawanRepo  repository.KaryawanRepository
	pelangganRepo repository.PelangganRepository
}

type UpdatePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func NewProfileController(
	uRepo repository.UserRepository, 
	aRepo repository.AdminRepository,
	kRepo repository.KaryawanRepository,
	pRepo repository.PelangganRepository,
) ProfileController {
	return &profileController{uRepo, aRepo, kRepo, pRepo}
}

func (ctrl *profileController) UpdateProfile(c *gin.Context) {
	// 1. Ambil ID User dari Token JWT (diset oleh middleware)
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak."})
		return
	}
	userID := uint(userIDFloat.(float64))

	// 2. Bind JSON dari Frontend
	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid."})
		return
	}

	// 3. Cari User untuk mengetahui RoleID-nya
	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan."})
		return
	}

	// 4. Update tabel 'users' (Username & Email)
	userUpdate := model.User{
		IDUser:   userID,
		Username: input.Username,
		Email:    input.Email,
	}
	if err := ctrl.userRepo.UpdateUser(&userUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update data akun."})
		return
	}

	// 5. Update tabel sesuai ROLE (Admin / Karyawan / Pelanggan)
	switch user.RoleID {
	case 1:
		admin := model.Admin{UserID: userID, NamaAdmin: input.Nama}
		ctrl.adminRepo.UpdateAdmin(&admin)
	case 2:
		karyawan := model.Karyawan{UserID: userID, NamaKaryawan: input.Nama}
		ctrl.karyawanRepo.UpdateKaryawan(&karyawan)
	case 3:
		pelanggan := model.Pelanggan{UserID: userID, NamaLengkap: input.Nama}
		ctrl.pelangganRepo.UpdatePelanggan(&pelanggan)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui!",
		"nama":    input.Nama,
	})
}

func (ctrl *profileController) UpdatePassword(c *gin.Context) {
	// Ambil ID dari Token
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak."})
		return
	}
	userID := uint(userIDFloat.(float64))

	var input UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap atau password kurang dari 6 karakter."})
		return
	}

	// Cari user di database
	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan."})
		return
	}

	// Cek apakah Old Password yang diinput cocok dengan di Database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password saat ini salah!"})
		return
	}

	// Hash Password Baru
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	
	// Update ke struct dan simpan
	user.Password = string(hashedPassword)
	if err := ctrl.userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui!"})
}