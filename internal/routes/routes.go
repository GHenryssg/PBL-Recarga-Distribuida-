package routes

import (
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/controllers"
	"github.com/gin-gonic/gin"
)

func ConfigurarRotas(router *gin.Engine) {
    router.POST("/reservePoint/:id", controllers.PostPoints)
    router.GET("/stations", controllers.GetAllRoutes)
}
