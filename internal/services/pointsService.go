package services

import (
    "errors"

    "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
    "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
)

func GetAllPoints() []models.PontoRecarga {
    return database.Pontos
}

func ReservePoints(ids []string) (reservados []string, indisponiveis []string, err error) {
    for _, id := range ids {
        reservado := false
        for i, ponto := range database.Pontos {
            if ponto.ID == id && ponto.Disponivel {
                database.Pontos[i].Disponivel = false
                reservados = append(reservados, id)
                reservado = true
                break
            }
        }
        if !reservado {
            indisponiveis = append(indisponiveis, id)
        }
    }

    if len(reservados) == 0 && len(indisponiveis) > 0 {
        err = errors.New("nenhum ponto foi reservado")
    }

    return reservados, indisponiveis, err
}