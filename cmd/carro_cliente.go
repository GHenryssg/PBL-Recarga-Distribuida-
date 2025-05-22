package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
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

func main() {
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

	// 4. Reservar os pontos via HTTP (requisição atômica, apenas exibe o resultado)
	idsParaReservar := []string{}
	for _, ponto := range pontosEscolhidos {
		idsParaReservar = append(idsParaReservar, ponto.ID)
	}
	fmt.Printf("\nTentando reservar os pontos: %v\n", idsParaReservar)
	urlReserva := fmt.Sprintf("%s/reserve-points/%s", servidor, strings.Join(idsParaReservar, ","))
	req, _ := http.NewRequest("POST", urlReserva, nil)
	respReserva, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Erro ao reservar pontos: %v\n", err)
		return
	}
	defer respReserva.Body.Close()
	var resultado map[string]interface{}
	json.NewDecoder(respReserva.Body).Decode(&resultado)
	fmt.Println("\nResultado da reserva:")
	if reservados, ok := resultado["reservados"]; ok {
		fmt.Printf("  Reservados: %v\n", reservados)
	}
	if indisponiveis, ok := resultado["indisponiveis"]; ok {
		fmt.Printf("  Indisponíveis: %v\n", indisponiveis)
	}
	if erro, ok := resultado["erro"]; ok {
		fmt.Printf("  Erro: %v\n", erro)
	}
	if _, ok := resultado["reservados"]; !ok {
		fmt.Println("Nenhum ponto foi reservado. Encerrando fluxo.")
		return
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

	// 6. Liberar os pontos (cancelar reservas) usando o novo endpoint sequencial
	idsParaCancelar := make([]string, len(pontosEscolhidos))
	for i, ponto := range pontosEscolhidos {
		idsParaCancelar[i] = ponto.ID
	}
	urlCancel := fmt.Sprintf("%s/cancel-reservation/%s", servidor, strings.Join(idsParaCancelar, ","))
	reqCancel, _ := http.NewRequest("POST", urlCancel, nil)
	respCancel, err := http.DefaultClient.Do(reqCancel)
	if err != nil {
		panic(err)
	}
	defer respCancel.Body.Close()
	var cancelResult map[string]interface{}
	json.NewDecoder(respCancel.Body).Decode(&cancelResult)
	fmt.Printf("Pontos liberados: %v\n", cancelResult)
}
