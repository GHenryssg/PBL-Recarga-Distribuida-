package main

import (
	mqtt "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/client"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.CarregarVariaveis()

	// Inicializa o servidor MQTT para comunicação assíncrona (cliente-servidor e entre empresas)
	// O backend escuta tópicos de reserva/cancelamento e responde via MQTT
	go mqtt.StartMQTT(config.MQTTBrokerURL)

	router := gin.Default()

	// Configura as rotas HTTP REST para comunicação síncrona
	router.SetTrustedProxies(nil) // Desativa o aviso sobre proxies confiáveis

	routes.ConfigurarRotas(router)

	// Inicializa o backend na porta configurada
	router.Run(":" + config.Porta)
}
