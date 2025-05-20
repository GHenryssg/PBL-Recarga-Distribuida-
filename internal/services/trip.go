package services

import (
	mqtt "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/client"
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
		liberado := false
		for i := range database.Pontos {
			if database.Pontos[i].ID == id {
				// Se for ponto local, libera normalmente
				if isPontoDaEmpresa(database.Pontos[i]) {
					database.Pontos[i].Disponivel = true
					AtualizarDisponibilidadeNasRotas(id, true)
					liberado = true
					break
				} else {
					// Se for remoto, envia mensagem MQTT para liberar
					nomeEmpresa := getEmpresaNomeDoID(database.Pontos[i].EmpresaID)
					mqtt.SolicitarLiberacaoRemota(id, nomeEmpresa)
					liberado = true
					break
				}
			}
		}
		if !liberado {
			println("[DEBUG] CancelarReservas: ponto", id, "não encontrado para liberar")
		}
	}
}
