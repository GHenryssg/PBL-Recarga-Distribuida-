package controllers

import (
    "net/http"

    "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
    "github.com/gin-gonic/gin"
)

func GetAllRoutes(c *gin.Context) {
    routes := services.GetAllRoutes()
    c.JSON(http.StatusOK, routes)
}

func GetRouteByID(c *gin.Context) {
    id := c.Param("id")
    route, err := services.GetRouteByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
        return
    }
    c.JSON(http.StatusOK, route)
}