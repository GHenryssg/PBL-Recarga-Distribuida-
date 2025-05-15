package controllers

import (
	"net/http"
	"strings"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func GetAllPoints(c *gin.Context) {
	points := services.GetAllPoints()
	c.JSON(http.StatusOK, points)
}

func PostPoints(c *gin.Context) {
	idsParam := c.Param("ids")
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "IDs não fornecidos na URL"})
		return
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Nenhum ID válido fornecido"})
		return
	}

	// Configurar o cliente MQTT corretamente
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt_broker:1883")
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao conectar ao MQTT"})
		return
	}
	defer mqttClient.Disconnect(250)

	reservados, indisponiveis, err := services.ReservePoints(ids, mqttClient)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservados":    reservados,
		"indisponiveis": indisponiveis,
	})
}
