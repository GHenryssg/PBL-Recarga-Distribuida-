package routes

import (
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/controllers"
	"github.com/gin-gonic/gin"
)

func ConfigurarRotas(router *gin.Engine) {
	router.GET("/points", controllers.GetAllPoints)                        // lista todos os pontos de recarga
	router.POST("/reserve-points/:ids", controllers.PostPoints)            // reserva múltiplos pontos (por IDs)
	router.GET("/routes", controllers.GetAllRoutes)                        // lista todas as rotas cadastradas
	router.GET("/routes/:id", controllers.GetRouteByID)                    // obtém uma rota específica por ID
	router.POST("/plan-trip", controllers.PlanTrip)                        // retorna pontos disponíveis entre origem e destino
	router.POST("/reserve-sequence", controllers.ReserveSequence)          // reserva uma sequência de pontos de recarga
	router.POST("/cancel-reservation", controllers.CancelReservation)      // cancela reservas por ID
	router.POST("/cancel-reservation/:ids", controllers.CancelPointsByIDs) // cancela reservas por lista de IDs na URL
	router.GET("/health", controllers.HealthCheck)                         // status do servidor
}
