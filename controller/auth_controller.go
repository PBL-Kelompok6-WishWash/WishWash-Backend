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
}

func NewAuthController(userRepo repository.UserRepository) AuthController {
	return &authController{userRepo}
}

// 4. Logika Register
func (ctrl *authController) Register(c *gin.Context) {
	var input RegisterInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// Mapping yang benar: hanya masukkan data yang benar-benar ada di struct model.User
	user := model.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
		RoleID:   input.RoleID, // Gunakan RoleID (uint)
	}

	if err := ctrl.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username atau Email mungkin sudah terdaftar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi akun berhasil!", "username": user.Username})
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