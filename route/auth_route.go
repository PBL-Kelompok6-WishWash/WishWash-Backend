package route // <-- Pastikan baris 1 adalah ini

import (
    "github.com/PBL-Kelompok6-WishWash/backend/controller"
    "github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authController controller.AuthController) {
    api := router.Group("/api/v1")
    {
        api.POST("/register", authController.Register)
        api.POST("/login", authController.Login)
    }
}