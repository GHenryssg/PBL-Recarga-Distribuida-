package services

import (
	"errors"

	mqtt "github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/client"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
)

func GetAllPoints() []models.PontoRecarga {
	return database.Pontos
}

func isPontoDaEmpresa(ponto models.PontoRecarga) bool {
	idEsperado := config.EmpresaNomeParaID[config.NomeEmpresa]
	// Comparação robusta: aceita tanto nome quanto ID numérico
	match := ponto.EmpresaID == config.NomeEmpresa || ponto.EmpresaID == idEsperado
	if !match {
		println("[DEBUG] isPontoDaEmpresa: ponto.EmpresaID=", ponto.EmpresaID, "config.NomeEmpresa=", config.NomeEmpresa, "idEsperado=", idEsperado, "-> NÃO É DA EMPRESA LOCAL")
	} else {
		println("[DEBUG] isPontoDaEmpresa: ponto.EmpresaID=", ponto.EmpresaID, "config.NomeEmpresa=", config.NomeEmpresa, "idEsperado=", idEsperado, "-> É DA EMPRESA LOCAL")
	}
	return match
}

// Atualiza o campo Disponivel nas rotas sempre que um ponto é reservado ou liberado
func AtualizarDisponibilidadeNasRotas(pontoID string, disponivel bool) {
	for r := range database.Rotas {
		for p := range database.Rotas[r].Pontos {
			if database.Rotas[r].Pontos[p].ID == pontoID {
				database.Rotas[r].Pontos[p].Disponivel = disponivel
			}
		}
	}
}

// Função auxiliar para buscar o nome da empresa a partir do ID
func getEmpresaNomeDoID(id string) string {
	for nome, num := range config.EmpresaNomeParaID {
		if num == id {
			return nome
		}
	}
	return id // fallback
}

func ReservePoints(ids []string) (reservados []string, indisponiveis []string, err error) {
	type pontoStatus struct {
		id         string
		local      bool
		disponivel bool
	}
	var statusList []pontoStatus

	// 1. Verifica disponibilidade de todos os pontos (NÃO reserva ainda)
	for _, id := range ids {
		found := false
		for _, ponto := range database.Pontos {
			if ponto.ID == id {
				found = true
				if isPontoDaEmpresa(ponto) {
					if ponto.Disponivel {
						statusList = append(statusList, pontoStatus{id, true, true})
						println("[DEBUG] Ponto", id, "(local)", "DISPONÍVEL")
					} else {
						statusList = append(statusList, pontoStatus{id, true, false})
						println("[DEBUG] Ponto", id, "(local)", "INDISPONÍVEL")
					}
				} else {
					// Verifica disponibilidade remota via MQTT
					nomeEmpresa := getEmpresaNomeDoID(ponto.EmpresaID)
					println("[DEBUG] Testando disponibilidade remota ponto:", id, "empresa:", nomeEmpresa, "EmpresaID:", ponto.EmpresaID)
					ok := mqtt.TestarDisponibilidadeRemota(id, nomeEmpresa)
					statusList = append(statusList, pontoStatus{id, false, ok})
					if ok {
						println("[DEBUG] Ponto", id, "(remoto)", "DISPONÍVEL")
					} else {
						println("[DEBUG] Ponto", id, "(remoto)", "INDISPONÍVEL ou FALHA MQTT")
					}
				}
				break
			}
		}
		if !found {
			statusList = append(statusList, pontoStatus{id, false, false})
			println("[DEBUG] Ponto", id, "NÃO ENCONTRADO")
		}
	}

	// 2. Se algum indisponível, retorna erro
	for _, st := range statusList {
		if !st.disponivel {
			println("[DEBUG] Algum ponto está indisponível. Nenhum será reservado.")
			for _, s := range statusList {
				if !s.disponivel {
					indisponiveis = append(indisponiveis, s.id)
					println("[DEBUG] Indisponível:", s.id)
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
					AtualizarDisponibilidadeNasRotas(st.id, false)
					reservados = append(reservados, st.id)
					println("[DEBUG] Reservado local:", st.id)
					break
				}
			}
		} else {
			nomeEmpresa := getEmpresaNomeDoID(getEmpresaIDDoPonto(st.id))
			println("[DEBUG] Solicitando reserva remota ponto:", st.id, "empresa:", nomeEmpresa)
			ok := mqtt.SolicitarReservaRemota(st.id, nomeEmpresa)
			if ok {
				AtualizarDisponibilidadeNasRotas(st.id, false)
				reservados = append(reservados, st.id)
				println("[DEBUG] Reservado remoto:", st.id)
			} else {
				// Rollback: desfaz reservas locais já feitas
				println("[DEBUG] Falha ao reservar remoto:", st.id, "- iniciando rollback")
				for _, rid := range reservados {
					for i, ponto := range database.Pontos {
						if ponto.ID == rid {
							database.Pontos[i].Disponivel = true
							AtualizarDisponibilidadeNasRotas(rid, true)
							println("[DEBUG] Rollback liberado:", rid)
						}
					}
				}
				indisponiveis = append(indisponiveis, st.id)
				err = errors.New("nenhum ponto foi reservado")
				return nil, indisponiveis, err
			}
		}
	}
	println("[DEBUG] Todos os pontos reservados com sucesso:", reservados)
	return reservados, nil, nil
}

// Função auxiliar para buscar o EmpresaID de um ponto
func getEmpresaIDDoPonto(pontoID string) string {
	for _, ponto := range database.Pontos {
		if ponto.ID == pontoID {
			return ponto.EmpresaID
		}
	}
	return ""
}
