package services

import (
	"fmt"
	"net/http"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
)

func BuscarPontosNaRota(origem, destino string) []models.PontoRecarga {
	var pontosDisponiveis []models.PontoRecarga

	for _, rota := range database.Rotas {
		if contemLocal(rota.Pontos, origem) && contemLocal(rota.Pontos, destino) {
			for _, ponto := range rota.Pontos {
				if ponto.Disponivel {
					pontosDisponiveis = append(pontosDisponiveis, ponto)
				}
			}
		}
	}

	return pontosDisponiveis
}

func contemLocal(pontos []models.PontoRecarga, local string) bool {
	for _, p := range pontos {
		if p.Localizacao == local {
			return true
		}
	}
	return false
}

func CancelarReservas(ids []string) {
	for _, id := range ids {
		for i := range database.Pontos {
			if database.Pontos[i].ID == id {
				if isPontoDaEmpresa(database.Pontos[i]) {
					database.Pontos[i].Disponivel = true
					AtualizarDisponibilidadeNasRotas(id, true)
					break
				} else {
					// Se for remoto, faz requisição HTTP para liberar via endpoint /cancel-reservation/:id
					empresaURL := getEmpresaURLDoPonto(database.Pontos[i].EmpresaID)
					if empresaURL != "" {
						urlCancel := fmt.Sprintf("%s/cancel-reservation/%s", empresaURL, id)
						req, _ := http.NewRequest("POST", urlCancel, nil)
						resp, err := http.DefaultClient.Do(req)
						if err == nil && resp.StatusCode == 200 {
							// Atualiza localmente também para refletir o estado
							database.Pontos[i].Disponivel = true
							AtualizarDisponibilidadeNasRotas(id, true)
						}
						if resp != nil {
							resp.Body.Close()
						}
					}
					break
				}
			}
		}
	}
}
