package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/config"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Inicializar o cliente MQTT uma vez e mantê-lo ativo
var mqttClient mqtt.Client
var once sync.Once

func getMQTTClient() mqtt.Client {
	once.Do(func() {
		opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt_broker:1883")
		opts.SetClientID("servidor-mqtt-" + config.NomeEmpresa)
		opts.SetAutoReconnect(true)
		opts.SetKeepAlive(120 * time.Second)
		mqttClient = mqtt.NewClient(opts)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	})
	// Verificar se o cliente MQTT já está conectado antes de reutilizá-lo
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Erro ao conectar ao MQTT:", token.Error())
	}
	return mqttClient
}

func GetAllPoints() []models.PontoRecarga {
	return database.Pontos
}

func ReservePoints(ids []string, mqttClient mqtt.Client) (reservados []string, indisponiveis []string, err error) {
	// Usar o cliente MQTT ativo
	mqttClient = getMQTTClient()

	for _, id := range ids {
		fmt.Printf("🔍 Verificando disponibilidade do ponto: %s\n", id)

		// Verificar disponibilidade via MQTT
		topic := "pontos/verificar"
		responseTopic := "pontos/resposta/" + id

		payload, _ := json.Marshal(id)
		token := mqttClient.Publish(topic, 0, false, payload)
		token.Wait()

		// Aguardar resposta
		var disponivel bool
		done := make(chan bool)
		mqttClient.Subscribe(responseTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
			var resp map[string]interface{}
			if err := json.Unmarshal(msg.Payload(), &resp); err == nil {
				disponivel = resp["disponivel"].(bool)
				if disponivel {
					done <- true
				}
			}
		})
		defer close(done)

		// Aumentar o timeout para 10 segundos
		select {
		case <-done:
			if !disponivel {
				indisponiveis = append(indisponiveis, id)
				continue
			}
		case <-time.After(10 * time.Second):
			indisponiveis = append(indisponiveis, id)
			continue
		}

		// Reservar ponto localmente
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
