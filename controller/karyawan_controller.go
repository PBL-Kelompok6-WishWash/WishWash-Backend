package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
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
		FotoKaryawan:       "",
		PlatNomor:          input.PlatNomor,
		JenisKendaraan:     input.JenisKendaraan,
		StatusKetersediaan: input.StatusKetersediaan,
	}

	if err := ctrl.karyawanRepo.CreateKaryawan(&karyawan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil karyawan"})
		return
	}

	if input.FotoKaryawan != "" {
		entityFolder := utils.BuildEntityFolder(karyawan.IDKaryawan, input.NamaKaryawan)
		fotoPath, err := utils.SaveBase64Image(input.FotoKaryawan, "karyawan", entityFolder, "profile_karyawan_"+input.NamaKaryawan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format foto karyawan tidak valid atau gagal disimpan"})
			return
		}
		karyawan.FotoKaryawan = fotoPath
		_ = ctrl.karyawanRepo.Update(&karyawan)
	}

	// 4. Ambil data lengkap (untuk preload User & Role)
	fullData, _ := ctrl.karyawanRepo.FindByID(karyawan.IDKaryawan)

	c.JSON(http.StatusCreated, gin.H{"message": "Karyawan berhasil ditambahkan!", "data": fullData})
}

type UpdateKaryawanInput struct {
	Username           *string `json:"username"`
	Email              *string `json:"email"`
	Password           *string `json:"password"`
	NamaKaryawan       *string `json:"nama_karyawan"`
	NoTelp             *string `json:"no_telp"`
	FotoKaryawan       *string `json:"foto_karyawan"`
	PlatNomor          *string `json:"plat_nomor"`
	JenisKendaraan     *string `json:"jenis_kendaraan"`
	StatusKetersediaan *string `json:"status_ketersediaan"`
}

func (ctrl *karyawanController) Update(c *gin.Context) {
	id := parseKaryawanID(c.Param("id"))
	var input UpdateKaryawanInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	karyawan, err := ctrl.karyawanRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan"})
		return
	}

	// Update User
	userUpdated := false
	if input.Username != nil {
		karyawan.User.Username = *input.Username
		userUpdated = true
	}
	if input.Email != nil {
		karyawan.User.Email = *input.Email
		userUpdated = true
	}
	if input.Password != nil && *input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		karyawan.User.Password = string(hashedPassword)
		userUpdated = true
	}

	if userUpdated {
		if err := ctrl.userRepo.UpdateUser(&karyawan.User); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate akun user"})
			return
		}
	}

	// Update Karyawan
	namaKaryawan := karyawan.NamaKaryawan
	if input.NamaKaryawan != nil {
		namaKaryawan = *input.NamaKaryawan
		karyawan.NamaKaryawan = *input.NamaKaryawan
	}
	if input.NoTelp != nil {
		karyawan.NoTelp = *input.NoTelp
	}
	if input.FotoKaryawan != nil {
		oldFotoPath := karyawan.FotoKaryawan
		entityFolder := utils.BuildEntityFolder(karyawan.IDKaryawan, namaKaryawan)

		if *input.FotoKaryawan == "" {
			// User klik X lalu simpan (hapus foto)
			karyawan.FotoKaryawan = ""
			if oldFotoPath != "" {
				utils.DeleteImageFile(oldFotoPath)
			}
		} else {
			// User upload foto baru
			fotoPath, err := utils.SaveBase64Image(*input.FotoKaryawan, "karyawan", entityFolder, "profile_karyawan_"+namaKaryawan)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan foto baru"})
				return
			}
			karyawan.FotoKaryawan = fotoPath

			// Hapus file lama jika path-nya berubah
			if oldFotoPath != "" && oldFotoPath != fotoPath {
				utils.DeleteImageFile(oldFotoPath)
			}
		}
	}
	if input.PlatNomor != nil {
		karyawan.PlatNomor = *input.PlatNomor
	}
	if input.JenisKendaraan != nil {
		karyawan.JenisKendaraan = *input.JenisKendaraan
	}
	if input.StatusKetersediaan != nil {
		karyawan.StatusKetersediaan = *input.StatusKetersediaan
	}

	if err := ctrl.karyawanRepo.Update(karyawan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate profil karyawan"})
		return
	}

	fullData, _ := ctrl.karyawanRepo.FindByID(id)
	c.JSON(http.StatusOK, gin.H{"message": "Data karyawan berhasil diperbarui!", "data": fullData})
}

func (ctrl *karyawanController) Delete(c *gin.Context) {
	id := parseKaryawanID(c.Param("id"))

	karyawan, err := ctrl.karyawanRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Karyawan tidak ditemukan"})
		return
	}

	// Hapus seluruh folder entity dari disk (bersih sekaligus)
	entityFolder := utils.BuildEntityFolder(karyawan.IDKaryawan, karyawan.NamaKaryawan)
	utils.DeleteImageFolder("karyawan", entityFolder)

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
