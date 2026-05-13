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

type KaryawanInput struct {
	Username           string `json:"username" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	Password           string `json:"password"`
	NamaKaryawan       string `json:"nama_karyawan" binding:"required"`
	NoTelp             string `json:"no_telp"`
	FotoKaryawan       string `json:"foto_karyawan"`
	PlatNomor          string `json:"plat_nomor"`
	JenisKendaraan     string `json:"jenis_kendaraan"`
	StatusKetersediaan string `json:"status_ketersediaan"`
}

type KaryawanController interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type karyawanController struct {
	karyawanRepo repository.KaryawanRepository
	userRepo     repository.UserRepository
}

func NewKaryawanController(kRepo repository.KaryawanRepository, uRepo repository.UserRepository) KaryawanController {
	return &karyawanController{kRepo, uRepo}
}

// Fungsi helper ada di pelanggan_controller.go, tapi kita bikin khusus aja kalau butuh
func parseKaryawanID(idStr string) uint {
	id, _ := strconv.ParseUint(idStr, 10, 32)
	return uint(id)
}

func (ctrl *karyawanController) GetAll(c *gin.Context) {
	karyawans, err := ctrl.karyawanRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data karyawan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": karyawans})
}

func (ctrl *karyawanController) GetByID(c *gin.Context) {
	id := parseKaryawanID(c.Param("id"))

	karyawan, err := ctrl.karyawanRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": karyawan})
}

func (ctrl *karyawanController) Create(c *gin.Context) {
	var input KaryawanInput

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
				case "NamaKaryawan":
					errorMsgs = append(errorMsgs, "Nama Karyawan tidak boleh kosong")
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMsgs, ", ")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data JSON tidak valid"})
		return
	}

	if input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak boleh kosong"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   2, // 2 = Role Karyawan
	}

	if err := ctrl.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun user"})
		return
	}

	karyawan := model.Karyawan{
		UserID:             user.IDUser,
		NamaKaryawan:       input.NamaKaryawan,
		NoTelp:             input.NoTelp,
		FotoKaryawan:       input.FotoKaryawan,
		PlatNomor:          input.PlatNomor,
		JenisKendaraan:     input.JenisKendaraan,
		StatusKetersediaan: input.StatusKetersediaan,
	}

	if err := ctrl.karyawanRepo.CreateKaryawan(&karyawan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil karyawan"})
		return
	}

	// 4. Ambil data lengkap (untuk preload User & Role)
	fullData, _ := ctrl.karyawanRepo.FindByID(karyawan.IDKaryawan)

	c.JSON(http.StatusCreated, gin.H{"message": "Karyawan berhasil ditambahkan!", "data": fullData})
}

func (ctrl *karyawanController) Update(c *gin.Context) {
	id := parseKaryawanID(c.Param("id"))
	var input KaryawanInput

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
				case "NamaKaryawan":
					errorMsgs = append(errorMsgs, "Nama Karyawan tidak boleh kosong")
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMsgs, ", ")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	karyawan, err := ctrl.karyawanRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan"})
		return
	}

	karyawan.User.Username = input.Username
	karyawan.User.Email = input.Email

	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		karyawan.User.Password = string(hashedPassword)
	}

	if err := ctrl.userRepo.UpdateUser(&karyawan.User); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate akun user"})
		return
	}

	karyawan.NamaKaryawan = input.NamaKaryawan
	karyawan.NoTelp = input.NoTelp
	karyawan.FotoKaryawan = input.FotoKaryawan
	karyawan.PlatNomor = input.PlatNomor
	karyawan.JenisKendaraan = input.JenisKendaraan
	karyawan.StatusKetersediaan = input.StatusKetersediaan

	if err := ctrl.karyawanRepo.Update(karyawan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate profil karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data karyawan berhasil diperbarui!"})
}

func (ctrl *karyawanController) Delete(c *gin.Context) {
	id := parseKaryawanID(c.Param("id"))

	karyawan, err := ctrl.karyawanRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan"})
		return
	}

	if err := ctrl.karyawanRepo.Delete(karyawan.IDKaryawan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus profil karyawan"})
		return
	}

	if err := ctrl.userRepo.DeleteUser(karyawan.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus akun user karyawan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Karyawan berhasil dihapus"})
}
