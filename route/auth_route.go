package route 

import (
    "github.com/PBL-Kelompok6-WishWash/backend/controller"
    "github.com/gin-gonic/gin"
)

func SetupAuthRoutes(apiGroup *gin.RouterGroup, authController controller.AuthController) {
    // Bikin sub-grup /auth
    auth := apiGroup.Group("/auth")
    {
        auth.POST("/register", authController.Register)
        auth.POST("/login", authController.Login)
    }
}