package main

import (
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/mqtt"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.CarregarVariaveis()

	// Use a variável global diretamente
	go mqtt.StartMQTT(config.MQTTBrokerURL)

	// Iniciar o serviço de resposta de pontos MQTT
	go mqtt.StartPointResponder(config.MQTTBrokerURL)

	router := gin.Default()

	// Configurar proxies confiáveis ou desativar o aviso
	router.SetTrustedProxies(nil) // Desativa o aviso sobre proxies confiáveis

	routes.ConfigurarRotas(router)

	// Mesma coisa aqui: usar config.Porta diretamente
	router.Run(":" + config.Porta)
}
