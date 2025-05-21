package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func StartMQTT(brokerURL string) {
	fmt.Printf("Iniciando cliente MQTT com ClientID: servidor-mqtt-%s\n", config.NomeEmpresa)
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("servidor-mqtt-" + config.NomeEmpresa)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient = client

	client.Subscribe("veiculos/solicitacao", 0, func(client mqtt.Client, msg mqtt.Message) {
		var req models.RequisicaoVeiculo
		if err := json.Unmarshal(msg.Payload(), &req); err != nil {
			fmt.Println("Erro ao decodificar MQTT:", err)
			return
		}
		fmt.Printf("Mensagem MQTT recebida: %+v\n", req)

		// Buscar pontos disponíveis na rota (sem ciclo de importação)
		pontosDisponiveis := []string{}
		for _, rota := range database.Rotas {
			if contemLocal(rota.Pontos, req.Local) && contemLocal(rota.Pontos, req.Destino) {
				for _, ponto := range rota.Pontos {
					if ponto.Disponivel {
						pontosDisponiveis = append(pontosDisponiveis, ponto.ID)
					}
				}
			}
		}

		resp := struct {
			VeiculoID         string   `json:"veiculo_id"`
			PontosDisponiveis []string `json:"pontos_disponiveis"`
		}{
			VeiculoID:         req.VeiculoID,
			PontosDisponiveis: pontosDisponiveis,
		}

		payload, _ := json.Marshal(resp)
		topic := "veiculos/resposta/" + req.VeiculoID
		client.Publish(topic, 0, false, payload)
		fmt.Printf("Respondido para %s\n", topic)
	})

	// Handler para reservas remotas de pontos
	topicReserva := "empresa/" + config.NomeEmpresa + "/reserva"
	fmt.Println("[MQTT] Subscribed to:", topicReserva)
	client.Subscribe(topicReserva, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println("[MQTT] Mensagem de reserva recebida neste container!", config.NomeEmpresa)
		var req struct {
			PontoID   string `json:"ponto_id"`
			EmpresaID string `json:"empresa_id"`
			Acao      string `json:"acao"`
			Origem    string `json:"origem"`
		}
		if err := json.Unmarshal(msg.Payload(), &req); err != nil {
			fmt.Println("Erro ao decodificar reserva remota:", err)
			return
		}
		if req.Acao == "liberar" {
			fmt.Printf("[MQTT] Pedido de liberação recebido para ponto %s neste container (%s)\n", req.PontoID, config.NomeEmpresa)
			liberado := false
			for i, ponto := range database.Pontos {
				if ponto.ID == req.PontoID && (ponto.EmpresaID == config.NomeEmpresa || ponto.EmpresaID == config.EmpresaNomeParaID[config.NomeEmpresa]) {
					database.Pontos[i].Disponivel = true
					for r := range database.Rotas {
						for p := range database.Rotas[r].Pontos {
							if database.Rotas[r].Pontos[p].ID == req.PontoID {
								database.Rotas[r].Pontos[p].Disponivel = true
							}
						}
					}
					fmt.Printf("[MQTT] Ponto %s liberado remotamente e marcado como disponível nas rotas.\n", req.PontoID)
					liberado = true
					break
				}
			}
			if !liberado {
				fmt.Printf("[MQTT] Liberação: ponto %s NÃO encontrado ou NÃO pertence a esta empresa (%s)\n", req.PontoID, config.NomeEmpresa)
			}
			return
		}
		for i, ponto := range database.Pontos {
			fmt.Printf("[MQTT] Verificando ponto %s: EmpresaID=%s NomeEmpresa=%s Disponivel=%v\n", ponto.ID, ponto.EmpresaID, config.NomeEmpresa, ponto.Disponivel)
			if ponto.ID == req.PontoID && (ponto.EmpresaID == config.NomeEmpresa || ponto.EmpresaID == config.EmpresaNomeParaID[config.NomeEmpresa]) && ponto.Disponivel {
				database.Pontos[i].Disponivel = false
				for r := range database.Rotas {
					for p := range database.Rotas[r].Pontos {
						if database.Rotas[r].Pontos[p].ID == req.PontoID {
							database.Rotas[r].Pontos[p].Disponivel = false
						}
					}
				}
				fmt.Printf("[MQTT] Ponto %s reservado remotamente e marcado como indisponível nas rotas.\n", req.PontoID)
				// Envia resposta positiva imediatamente e retorna
				resp := struct {
					PontoID    string `json:"ponto_id"`
					EmpresaID  string `json:"empresa_id"`
					Disponivel bool   `json:"disponivel"`
					Reservado  bool   `json:"reservado"`
					Mensagem   string `json:"mensagem"`
				}{
					PontoID:    req.PontoID,
					EmpresaID:  config.NomeEmpresa,
					Disponivel: false,
					Reservado:  true,
					Mensagem:   "",
				}
				payload, _ := json.Marshal(resp)
				respTopic := "empresa/" + req.Origem + "/resposta/" + req.PontoID
				client.Publish(respTopic, 0, false, payload)
				return
			}
		}
		// Se não encontrou ou não conseguiu reservar, responde negativo
		resp := struct {
			PontoID    string `json:"ponto_id"`
			EmpresaID  string `json:"empresa_id"`
			Disponivel bool   `json:"disponivel"`
			Reservado  bool   `json:"reservado"`
			Mensagem   string `json:"mensagem"`
		}{
			PontoID:    req.PontoID,
			EmpresaID:  config.NomeEmpresa,
			Disponivel: true,
			Reservado:  false,
			Mensagem:   "Ponto indisponível ou não pertence a esta empresa.",
		}
		payload, _ := json.Marshal(resp)
		respTopic := "empresa/" + req.Origem + "/resposta/" + req.PontoID
		client.Publish(respTopic, 0, false, payload)
	})
}

// Função auxiliar local (evita ciclo de importação)
func contemLocal(pontos []models.PontoRecarga, local string) bool {
	for _, p := range pontos {
		if p.Localizacao == local {
			return true
		}
	}
	return false
}

// Solicita reserva remota de um ponto para outra empresa via MQTT
func SolicitarReservaRemota(pontoID, empresaID string) bool {
	// Monta a requisição
	req := map[string]interface{}{
		"ponto_id":   pontoID,
		"empresa_id": empresaID,
		"acao":       "reservar",
		"origem":     config.NomeEmpresa,
	}
	payload, _ := json.Marshal(req)
	// Tópico de requisição e resposta
	reqTopic := "empresa/" + empresaID + "/reserva"
	respTopic := "empresa/" + config.NomeEmpresa + "/resposta/" + pontoID

	// Canal para resposta
	ch := make(chan bool, 1)

	// Handler de resposta
	token := mqttClient.Subscribe(respTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var resp models.ReservaPontoResponse
		err := json.Unmarshal(msg.Payload(), &resp)
		fmt.Printf("[MQTT][DEBUG] Resposta recebida para ponto %s: err=%v, resp=%+v\n", pontoID, err, resp)
		if err == nil && resp.PontoID == pontoID {
			fmt.Printf("[MQTT][DEBUG] Valor de resp.Reservado para ponto %s: %v\n", pontoID, resp.Reservado)
			ch <- resp.Reservado
		} else {
			ch <- false
		}
	})
	token.Wait()

	// Publica requisição
	mqttClient.Publish(reqTopic, 0, false, payload)

	// Aguarda resposta (timeout simples)
	select {
	case ok := <-ch:
		mqttClient.Unsubscribe(respTopic)
		// Pequeno delay para evitar race condition entre múltiplos subscribes
		time.Sleep(100 * time.Millisecond)
		return ok
	// Timeout de 3 segundos
	case <-time.After(3 * time.Second):
		mqttClient.Unsubscribe(respTopic)
		return false
	}
}

// Solicita liberação remota de um ponto para outra empresa via MQTT
func SolicitarLiberacaoRemota(pontoID, empresaID string) bool {
	req := map[string]interface{}{
		"ponto_id":   pontoID,
		"empresa_id": empresaID,
		"acao":       "liberar",
		"origem":     config.NomeEmpresa,
	}
	payload, _ := json.Marshal(req)
	reqTopic := "empresa/" + empresaID + "/reserva"
	respTopic := "empresa/" + config.NomeEmpresa + "/resposta/" + pontoID

	ch := make(chan bool, 1)
	token := mqttClient.Subscribe(respTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		ch <- true // Não precisa de resposta detalhada para liberação
	})
	token.Wait()
	mqttClient.Publish(reqTopic, 0, false, payload)
	select {
	case <-ch:
		mqttClient.Unsubscribe(respTopic)
		return true
	case <-time.After(2 * time.Second):
		mqttClient.Unsubscribe(respTopic)
		return false
	}
}

// Testa se um ponto remoto está disponível via MQTT
func TestarDisponibilidadeRemota(pontoID, empresaID string) bool {
	// Monta a requisição de verificação
	req := map[string]interface{}{
		"ponto_id":   pontoID,
		"empresa_id": empresaID,
		"acao":       "verificar",
		"origem":     config.NomeEmpresa,
	}
	payload, _ := json.Marshal(req)
	// Tópico de requisição e resposta
	reqTopic := "empresa/" + empresaID + "/reserva"
	respTopic := "empresa/" + config.NomeEmpresa + "/resposta/" + pontoID

	ch := make(chan bool, 1)

	token := mqttClient.Subscribe(respTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var resp models.ReservaPontoResponse
		if err := json.Unmarshal(msg.Payload(), &resp); err == nil && resp.PontoID == pontoID {
			ch <- resp.Disponivel
		} else {
			ch <- false
		}
	})
	token.Wait()

	mqttClient.Publish(reqTopic, 0, false, payload)

	select {
	case ok := <-ch:
		mqttClient.Unsubscribe(respTopic)
		return ok
	case <-time.After(3 * time.Second):
		mqttClient.Unsubscribe(respTopic)
		return false
	}
}
