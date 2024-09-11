package main

import (
	"xenon-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		panic(err.Error())
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	routes.UserRouter(router)
	router.Run()
}
