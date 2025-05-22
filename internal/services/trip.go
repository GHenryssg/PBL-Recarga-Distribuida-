package services

import (
	"fmt"
	"net/http"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
)

// BuscarPontosNaRota: Função de serviço chamada pelo controller de rotas para buscar pontos disponíveis
// Utilizada quando o cliente solicita rotas possíveis entre origem e destino
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

// contemLocal: Função auxiliar usada internamente para verificar se um ponto está em uma rota
func contemLocal(pontos []models.PontoRecarga, local string) bool {
	for _, p := range pontos {
		if p.Localizacao == local {
			return true
		}
	}
	return false
}

// CancelarReservas: Serviço chamado pelo controller de reservas e também via MQTT para cancelar reservas
// Se o ponto pertence a esta empresa, libera localmente; se for remoto, faz requisição HTTP para o backend da empresa dona do ponto
// Este método é fundamental para garantir rollback distribuído em caso de falha na reserva atômica
func CancelarReservas(ids []string) {
	for _, id := range ids {
		for i := range database.Pontos {
			if database.Pontos[i].ID == id {
				if isPontoDaEmpresa(database.Pontos[i]) {
					// Libera localmente
					database.Pontos[i].Disponivel = true
					AtualizarDisponibilidadeNasRotas(id, true)
					break
				} else {
					// Se for remoto, faz requisição para liberar via endpoint /cancel-reservation/:id
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
