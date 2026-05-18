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
	Username       string `json:"username"`
	Email          string `json:"email"`
	Nama           string `json:"nama"` // Menampung NamaAdmin / NamaKaryawan / NamaLengkap
	NoTelp         string `json:"no_telp"`
	FotoPelanggan  string `json:"foto_pelanggan"`
	PlatNomor      string `json:"plat_nomor"`
	JenisKendaraan string `json:"jenis_kendaraan"`
}

type ProfileController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type profileController struct {
	userRepo      repository.UserRepository
	adminRepo     repository.AdminRepository
	karyawanRepo  repository.KaryawanRepository
	pelangganRepo repository.PelangganRepository
	alamatRepo    repository.AlamatRepository
}

type UpdatePasswordInput struct {
	PasswordLama       string `json:"password_lama"`
	PasswordBaru       string `json:"password_baru"`
	KonfirmasiPassword string `json:"konfirmasi_password"`
	
	// Backward compatibility fields for Next.js web dashboard
	OldPassword     string `json:"oldPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

func NewProfileController(
	uRepo repository.UserRepository, 
	aRepo repository.AdminRepository,
	kRepo repository.KaryawanRepository,
	pRepo repository.PelangganRepository,
	alRepo repository.AlamatRepository,
) ProfileController {
	return &profileController{uRepo, aRepo, kRepo, pRepo, alRepo}
}

func (ctrl *profileController) GetProfile(c *gin.Context) {
	// 1. Ambil ID User & Role ID dari Token JWT
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak."})
		return
	}
	userID := uint(userIDFloat.(float64))

	roleIDFloat, existsRole := c.Get("id_role")
	if !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak."})
		return
	}
	roleID := int(roleIDFloat.(float64))

	// 2. Siapkan variabel penampung data
	var profileData interface{}
	var err error

	// 3. Tarik data dari Database sesuai Role-nya
	switch roleID {
	case 1:
		profileData, err = ctrl.adminRepo.FindByUserID(userID)
	case 2:
		profileData, err = ctrl.karyawanRepo.FindByUserID(userID)
	case 3:
		pelanggan, errP := ctrl.pelangganRepo.FindByUserID(userID)
		err = errP
		if errP == nil {
			var alamatLengkap string
			var tipeAlamat string
			alamat, errAlamat := ctrl.alamatRepo.FindByPelangganID(pelanggan.IDPelanggan)
			if errAlamat == nil && alamat != nil {
				alamatLengkap = alamat.AlamatLengkap
				tipeAlamat = alamat.TipeAlamat
			} else {
				alamatLengkap = "Alamat belum diatur"
				tipeAlamat = "Rumah"
			}
			profileData = gin.H{
				"pelanggan":      pelanggan,
				"alamat_lengkap": alamatLengkap,
				"tipe_alamat":    tipeAlamat,
			}
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tidak dikenali sistem."})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data profil tidak ditemukan di database."})
		return
	}

	// 4. Kembalikan data utuh
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Profil Berhasil Diambil",
		"data":    profileData,
	})
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
		existingAdmin, err := ctrl.adminRepo.FindByUserID(userID)
		if err == nil {
			if input.Nama != "" {
				existingAdmin.NamaAdmin = input.Nama
			}
			ctrl.adminRepo.UpdateAdmin(existingAdmin)
		}
	case 2:
		existingKaryawan, err := ctrl.karyawanRepo.FindByUserID(userID)
		if err == nil {
			if input.Nama != "" {
				existingKaryawan.NamaKaryawan = input.Nama
			}
			if input.NoTelp != "" {
				existingKaryawan.NoTelp = input.NoTelp
			}
			if input.FotoPelanggan != "" {
				existingKaryawan.FotoKaryawan = input.FotoPelanggan
			}
			if input.PlatNomor != "" {
				existingKaryawan.PlatNomor = input.PlatNomor
			}
			if input.JenisKendaraan != "" {
				existingKaryawan.JenisKendaraan = input.JenisKendaraan
			}
			ctrl.karyawanRepo.UpdateKaryawan(existingKaryawan)
		}
	case 3:
		existingPelanggan, err := ctrl.pelangganRepo.FindByUserID(userID)
		if err == nil {
			if input.Nama != "" {
				existingPelanggan.NamaLengkap = input.Nama
			}
			if input.NoTelp != "" {
				existingPelanggan.NoTelp = input.NoTelp
			}
			if input.FotoPelanggan != "" {
				existingPelanggan.FotoPelanggan = input.FotoPelanggan
			}
			ctrl.pelangganRepo.UpdatePelanggan(existingPelanggan)
		}
	}

	// 6. Tarik data profil terbaru untuk dikembalikan utuh ke client
	var profileData interface{}
	switch user.RoleID {
	case 1:
		profileData, _ = ctrl.adminRepo.FindByUserID(userID)
	case 2:
		profileData, _ = ctrl.karyawanRepo.FindByUserID(userID)
	case 3:
		pelanggan, errP := ctrl.pelangganRepo.FindByUserID(userID)
		if errP == nil {
			var alamatLengkap string
			var tipeAlamat string
			alamat, errAlamat := ctrl.alamatRepo.FindByPelangganID(pelanggan.IDPelanggan)
			if errAlamat == nil && alamat != nil {
				alamatLengkap = alamat.AlamatLengkap
				tipeAlamat = alamat.TipeAlamat
			} else {
				alamatLengkap = "Alamat belum diatur"
				tipeAlamat = "Rumah"
			}
			profileData = gin.H{
				"pelanggan":      pelanggan,
				"alamat_lengkap": alamatLengkap,
				"tipe_alamat":    tipeAlamat,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil berhasil diperbarui!",
		"data":    profileData,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid."})
		return
	}

	// 1. Ambil password lama dengan toleransi dual-language
	passwordLama := input.PasswordLama
	if passwordLama == "" {
		passwordLama = input.OldPassword
	}

	// 2. Ambil password baru dengan toleransi dual-language
	passwordBaru := input.PasswordBaru
	if passwordBaru == "" {
		passwordBaru = input.NewPassword
	}

	// 3. Ambil konfirmasi password dengan toleransi dual-language
	konfirmasiPassword := input.KonfirmasiPassword
	if konfirmasiPassword == "" {
		konfirmasiPassword = input.ConfirmPassword
	}

	// 4. Validasi keberadaan data wajib
	if passwordLama == "" || passwordBaru == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap! Harap masukkan password lama dan password baru."})
		return
	}

	// 5. Validasi panjang password baru
	if len(passwordBaru) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password baru harus minimal 6 karakter!"})
		return
	}

	// 6. Cek konfirmasi password (hanya jika dikirimkan oleh client/postman)
	if konfirmasiPassword != "" && passwordBaru != konfirmasiPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Konfirmasi password baru tidak cocok!"})
		return
	}

	// Cari user di database
	user, err := ctrl.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan."})
		return
	}

	// Cek apakah Password Lama cocok dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordLama))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password saat ini salah!"})
		return
	}

	// Hash Password Baru
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordBaru), bcrypt.DefaultCost)
	
	// Update ke struct dan simpan
	user.Password = string(hashedPassword)
	if err := ctrl.userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diperbarui!"})
}