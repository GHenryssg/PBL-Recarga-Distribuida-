package services

import (
    "errors"

    "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
    "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
)

func GetAllRoutes() []models.Rota {
    return database.Rotas
}

func GetRouteByID(id string) (models.Rota, error) {
    for _, rota := range database.Rotas {
        if rota.ID == id {
            return rota, nil
        }
    }
    return models.Rota{}, errors.New("rota não encontrada")
}