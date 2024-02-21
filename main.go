package main

import (

	"backend-assessment/routes"
	"backend-assessment/services"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	services.Initialize()
	routes.Initialize()

	r.Run()
}