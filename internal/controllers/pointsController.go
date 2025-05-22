package controllers

import (
	"net/http"
	"strings"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
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
	ids := []string{}
	for _, id := range strings.Split(idsParam, ",") {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			ids = append(ids, trimmed)
		}
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

// Endpoint para cancelar reservas por lista de IDs na URL
func CancelPointsByIDs(c *gin.Context) {
	idsParam := c.Param("ids")
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "IDs não informados"})
		return
	}
	ids := strings.Split(idsParam, ",")
	cancelados, naoCancelados, err := services.CancelPoints(ids)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error(), "nao_cancelados": naoCancelados})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cancelados":     cancelados,
		"nao_cancelados": naoCancelados,
	})
}
