package main

import (
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/mqtt"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	go mqtt.StartMQTT() // Inicia o cliente MQTT em paraleloo

	router := gin.Default()
	routes.ConfigurarRotas(router)
	router.Run(":8080")
}
