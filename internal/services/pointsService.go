package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/database"
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func GetAllPoints() []models.PontoRecarga {
    return database.Pontos
}

func ReservePoints(ids []string, mqttClient mqtt.Client) (reservados []string, indisponiveis []string, err error) {
    for _, id := range ids {
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
            }
            done <- true
        })

        select {
        case <-done:
            if !disponivel {
                indisponiveis = append(indisponiveis, id)
                continue
            }
        case <-time.After(5 * time.Second):
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