package services

import (
	"errors"
	"fmt"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/client"
)

func GetAllPoints() []models.PontoRecarga {
	return database.Pontos
}

func isPontoDaEmpresa(ponto models.PontoRecarga) bool {
	idEsperado := config.EmpresaNomeParaID[config.NomeEmpresa]
	return ponto.EmpresaID == config.NomeEmpresa || ponto.EmpresaID == idEsperado
}

func ReservePoints(ids []string) (reservados []string, indisponiveis []string, err error) {
	type pontoStatus struct {
		id         string
		local      bool
		disponivel bool
	}
	var statusList []pontoStatus

	// 1. Verifica disponibilidade de todos os pontos
	for _, id := range ids {
		found := false
		for i, ponto := range database.Pontos {
			if ponto.ID == id {
				fmt.Printf("ID: %s, EmpresaID do ponto: %s, NomeEmpresa config: %s\n", id, ponto.EmpresaID, config.NomeEmpresa)
				found = true
				if isPontoDaEmpresa(ponto) {
					// Reserva local
					if ponto.Disponivel {
						database.Pontos[i].Disponivel = false
						reservados = append(reservados, id)
					} else {
						statusList = append(statusList, pontoStatus{id, true, false})
					}
				} else {
					// Reserva remota via MQTT
					ok := mqtt.SolicitarReservaRemota(id, ponto.EmpresaID)
					if ok {
						reservados = append(reservados, id)
					} else {
						statusList = append(statusList, pontoStatus{id, false, false})
					}
				}
				break
			}
		}
		if !found {
			statusList = append(statusList, pontoStatus{id, false, false})
		}
	}

	// 2. Se algum indisponível, retorna erro
	for _, st := range statusList {
		if !st.disponivel {
			for _, s := range statusList {
				if !s.disponivel {
					indisponiveis = append(indisponiveis, s.id)
				}
			}
			err = errors.New("nenhum ponto foi reservado")
			return nil, indisponiveis, err
		}
	}

	// 3. Reserva todos (agora é garantido que todos estão disponíveis)
	for _, st := range statusList {
		if st.local {
			for i, ponto := range database.Pontos {
				if ponto.ID == st.id {
					database.Pontos[i].Disponivel = false
					reservados = append(reservados, st.id)
					break
				}
			}
		} else {
			ok := mqtt.SolicitarReservaRemota(st.id, "empresa_id_aqui") // Ajuste para pegar o EmpresaID correto
			if ok {
				reservados = append(reservados, st.id)
			} else {
				// Rollback: desfaz reservas locais já feitas
				for _, rid := range reservados {
					for i, ponto := range database.Pontos {
						if ponto.ID == rid {
							database.Pontos[i].Disponivel = true
						}
					}
				}
				indisponiveis = append(indisponiveis, st.id)
				err = errors.New("nenhum ponto foi reservado")
				return nil, indisponiveis, err
			}
		}
	}
	return reservados, nil, nil
}
