package controllers

import (
    "net/http"

    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
    "github.com/gin-gonic/gin"
)

func GetAllPoints(c *gin.Context) {
    points := services.GetAllPoints()
    c.JSON(http.StatusOK, points)
}

func PostPoints(c *gin.Context) {
    var ids []string
    if err := c.BindJSON(&ids); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"erro": "Formato de dados inválido"})
        return
    }

    // Configurar o cliente MQTT corretamente
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
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