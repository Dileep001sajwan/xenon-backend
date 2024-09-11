package routes

import (
	"xenon-backend/controller"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.POST("/user/signup", controller.SignUp())
	router.POST("/user/login", controller.Login())
}
