package main

import (
	"github.com/gin-gonic/gin"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/controllers"
)

func main() {
	router := gin.Default()

	// Rotas principais
	router.GET("/health", controllers.HealthCheck)

	// Pontos de recarga
	router.GET("/pontos", controllers.GetAllPoints)
	router.POST("/pontos/reservar", controllers.PostPoints)

	// Rotas
	router.GET("/rotas", controllers.GetAllRoutes)
	router.GET("/rotas/:id", controllers.GetRouteByID)

	// Viagem
	router.POST("/viagem/planejar", controllers.PlanTrip)
	router.POST("/viagem/reservar", controllers.ReserveSequence)
	router.POST("/viagem/cancelar", controllers.CancelReservation)

	// Executar servidor
	router.Run(":8080")
}
