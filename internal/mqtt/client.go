package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/services"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMQTT(brokerURL string) {
	fmt.Printf("Iniciando cliente MQTT com ClientID: servidor-mqtt-%s\n", config.NomeEmpresa)
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("servidor-mqtt-" + config.NomeEmpresa)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Subscribe("veiculos/solicitacao", 0, func(client mqtt.Client, msg mqtt.Message) {
		var req models.RequisicaoVeiculo
		if err := json.Unmarshal(msg.Payload(), &req); err != nil {
			fmt.Println("Erro ao decodificar MQTT:", err)
			return
		}

		pontos := services.BuscarPontosNaRota(req.Local, req.Destino)

		resp := models.RespostaServidor{
			VeiculoID:         req.VeiculoID,
			PontosDisponiveis: pontos,
		}

		payload, _ := json.Marshal(resp)
		topic := "veiculos/resposta/" + req.VeiculoID

		client.Publish(topic, 0, false, payload)

		fmt.Printf("📡 Respondido para %s\n", topic)
	})
}
