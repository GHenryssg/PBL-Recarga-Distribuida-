package mqtt

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
)

const (
	brokerURL       = "tcp://localhost:1883"
	topicSolicitacao = "veiculos/solicitacao"
	topicResposta    = "veiculos/resposta/"
)

func StartMQTT() {
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("servidor-rede")

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Escuta por requisições de viagem
	client.Subscribe(topicSolicitacao, 0, func(client mqtt.Client, msg mqtt.Message) {
		var req models.RequisicaoVeiculo
		if err := json.Unmarshal(msg.Payload(), &req); err != nil {
			fmt.Println("Erro ao decodificar mensagem:", err)
			return
		}

		// Processa a rota e encontra pontos
		pontos := services.BuscarPontosNaRota(req.Local, req.Destino)

		resp := models.RespostaServidor{
			VeiculoID:         req.VeiculoID,
			PontosDisponiveis: pontos,
		}

		payload, _ := json.Marshal(resp)
		topic := topicResposta + req.VeiculoID

		client.Publish(topic, 0, false, payload)

		fmt.Printf("📨 Requisição de %s processada. Enviada para %s.\n", req.VeiculoID, topic)
	})
}
