package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartPointResponder(brokerURL string) {
	fmt.Println("Iniciando serviço de resposta de pontos MQTT")
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	// Tornar o ClientID único usando o nome da empresa
	opts.SetClientID("point-responder-" + config.NomeEmpresa)

	// Garantir que o cliente MQTT permaneça conectado
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(120 * time.Second)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Subscribe("pontos/verificar", 0, func(client mqtt.Client, msg mqtt.Message) {
		var pointID string
		if err := json.Unmarshal(msg.Payload(), &pointID); err != nil {
			fmt.Println("Erro ao decodificar mensagem MQTT:", err)
			return
		}

		// Verificar se o ponto pertence à empresa antes de responder
		responsavel := false
		for _, empresa := range database.Empresas {
			if empresa.Nome == config.NomeEmpresa {
				for _, ponto := range empresa.Pontos {
					if ponto.ID == pointID {
						responsavel = true
						break
					}
				}
			}
		}
		if !responsavel {
			return
		}

		disponivel := false
		for i, ponto := range database.Pontos {
			if ponto.ID == pointID {
				// Centralizar a sincronização e evitar múltiplas respostas
				if ponto.Disponivel {
					database.Pontos[i].Disponivel = false
					disponivel = true
				} else {
					disponivel = false
				}
				break
			}
		}

		// Remover logs desnecessários
		response := map[string]interface{}{
			"disponivel": disponivel,
		}
		payload, _ := json.Marshal(response)
		responseTopic := "pontos/resposta/" + pointID
		client.Publish(responseTopic, 0, false, payload)
	})
}
