package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/eclipse/paho.mqtt.golang"
)

type PontoRecarga struct {
	ID          string `json:"id"`
	Localizacao string `json:"localizacao"`
	Disponivel  bool   `json:"disponivel"`
	EmpresaID   string `json:"empresa_id"`
}

type Rota struct {
	ID     string         `json:"id"`
	Nome   string         `json:"nome"`
	Pontos []PontoRecarga `json:"pontos"`
}

// Função auxiliar para converter um slice em JSON
func toJSON(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

var mqttClient mqtt.Client

func initMQTT() {
    opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883") // Substitua pelo endereço do seu broker MQTT
    opts.SetClientID("carro_cliente")
    mqttClient = mqtt.NewClient(opts)
    if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    fmt.Println("Conectado ao broker MQTT")
}


// Adicionar função para reservar ponto via MQTT
func reservarPontoViaMQTT(pontoID string) {
    topic := fmt.Sprintf("empresa/%s/reserva", pontoID)
    mensagem := fmt.Sprintf(`{"ponto_id": "%s", "acao": "reservar"}`, pontoID)
    token := mqttClient.Publish(topic, 0, false, mensagem)
    token.Wait()
    if token.Error() != nil {
        fmt.Printf("Erro ao publicar no tópico %s: %v\n", topic, token.Error())
    } else {
        fmt.Printf("Mensagem publicada no tópico %s: %s\n", topic, mensagem)
    }
}


func main() {
	rand.Seed(time.Now().UnixNano())

    // Inicializar o cliente MQTT
    initMQTT()
	
	rand.Seed(time.Now().UnixNano())
	servidor := os.Getenv("SERVER_URL")
	if servidor == "" {
		servidor = "http://localhost:8085"
	}

	// 1. Buscar rotas disponíveis
	resp, err := http.Get(servidor + "/routes")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var rotas []Rota
	if err := json.NewDecoder(resp.Body).Decode(&rotas); err != nil {
		panic(err)
	}
	fmt.Println("Rotas disponíveis:")
	for i, r := range rotas {
		fmt.Printf("%d - %s\n", i+1, r.Nome)
	}

	// 2. Escolhe uma rota
	rota := rotas[rand.Intn(len(rotas))]
	fmt.Printf("\nRota escolhida: %s\n", rota.Nome)
	fmt.Println("Pontos dessa rota:")
	for i, p := range rota.Pontos {
		fmt.Printf("%d - %s (%s) [Disponível: %v]\n", i+1, p.ID, p.Localizacao, p.Disponivel)
	}

	// 3. Escolher até 3 pontos disponíveis aleatórios
	var pontosDisponiveis []PontoRecarga
	for _, p := range rota.Pontos {
		if p.Disponivel {
			pontosDisponiveis = append(pontosDisponiveis, p)
		}
	}
	if len(pontosDisponiveis) == 0 {
		fmt.Println("Nenhum ponto disponível para reservar nesta rota!")
		return
	}

	// Selecionar até 3 pontos aleatórios
	numReservas := 3
	if len(pontosDisponiveis) < numReservas {
		numReservas = len(pontosDisponiveis)
	}
	pontosEscolhidos := pontosDisponiveis[:numReservas]
	fmt.Println("\nPontos escolhidos:")
	for _, ponto := range pontosEscolhidos {
		fmt.Printf("- %s (%s)\n", ponto.ID, ponto.Localizacao)
	}

	//Fazer verificação se os pontos escolhidos estão disponíveis e mostrar sua disponibilidade nas outras empresas

    // 4. Reservar os pontos via MQTT
    fmt.Println("\nReservando os pontos ...")
    for _, ponto := range pontosEscolhidos {
        reservarPontoViaMQTT(ponto.ID)
    }

	// 5. Simular viagem em partes
	fmt.Println("\nIniciando a viagem...")
	for i := 0; i < len(pontosEscolhidos); i++ {
		var origem, destino string
		if i == 0 {
			origem = "Início da rota"
		} else {
			origem = pontosEscolhidos[i-1].Localizacao
		}
		destino = pontosEscolhidos[i].Localizacao

		duracao := rand.Intn(3) + 7 // Tempo entre 3 e 7 segundos para cada trecho
		fmt.Printf("Viajando de %s para %s...\n", origem, destino)
		time.Sleep(time.Duration(duracao) * time.Second)
		fmt.Printf("Chegou em %s!\n", destino)
	}
	fmt.Println("Viagem concluída!")

	// 6. Liberar os pontos (cancelar reservas)
	urlCancel := fmt.Sprintf("%s/cancel-reservation", servidor)
	ids := make([]string, len(pontosEscolhidos))
	for i, ponto := range pontosEscolhidos {
		ids[i] = ponto.ID
	}
	body := fmt.Sprintf(`{"ids":%s}`, toJSON(ids))
	respCancel, err := http.Post(urlCancel, "application/json", strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer respCancel.Body.Close()
	var cancelResult map[string]interface{}
	json.NewDecoder(respCancel.Body).Decode(&cancelResult)
	fmt.Printf("Pontos liberados: %v\n", cancelResult)
}
