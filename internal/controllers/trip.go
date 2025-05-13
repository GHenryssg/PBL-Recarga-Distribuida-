package controllers

import (
	"net/http"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func PlanTrip(c *gin.Context) {
	var req models.PlanejamentoViagem
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	pontos := services.BuscarPontosNaRota(req.Origem, req.Destino)
	c.JSON(http.StatusOK, gin.H{"pontos_disponiveis": pontos})
}

func ReserveSequence(c *gin.Context) {
	var req models.SequenciaReserva
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato inválido"})
		return
	}

	ids := []string{}
	for _, ponto := range req.Pontos {
		ids = append(ids, ponto.ID)
	}

	// Configurar o cliente MQTT corretamente
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao conectar ao MQTT"})
		return
	}
	defer mqttClient.Disconnect(250)

	// Chamar a função ReservePoints com o cliente MQTT
	reservados, indisponiveis, err := services.ReservePoints(ids, mqttClient)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error(), "indisponiveis": indisponiveis})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservados":    reservados,
		"indisponiveis": indisponiveis,
	})
}

func CancelReservation(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato inválido"})
		return
	}

	services.CancelarReservas(req.IDs)
	c.JSON(http.StatusOK, gin.H{"mensagem": "Reservas canceladas"})
}
