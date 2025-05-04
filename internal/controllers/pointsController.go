package controllers

import (
    "net/http"

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

    reservados, indisponiveis, err := services.ReservePoints(ids)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "reservados":    reservados,
        "indisponiveis": indisponiveis,
    })
}