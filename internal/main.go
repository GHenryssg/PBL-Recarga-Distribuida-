package main

import (
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.ConfigurarRotas(router)
	router.Run(":8080")
}
