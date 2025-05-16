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

func main() {
	rand.Seed(time.Now().UnixNano())
	servidor := os.Getenv("SERVER_URL")
	if servidor == "" {
		servidor = "http://localhost:8081"
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

	// 2. Sortear uma rota
	rota := rotas[rand.Intn(len(rotas))]
	fmt.Printf("\nRota sorteada: %s\n", rota.Nome)
	fmt.Println("Pontos dessa rota:")
	for i, p := range rota.Pontos {
		fmt.Printf("%d - %s (%s) [Disponível: %v]\n", i+1, p.ID, p.Localizacao, p.Disponivel)
	}

	// 3. Escolher um ponto disponível aleatório
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
	pontoEscolhido := pontosDisponiveis[rand.Intn(len(pontosDisponiveis))]
	fmt.Printf("\nPonto escolhido: %s (%s)\n", pontoEscolhido.ID, pontoEscolhido.Localizacao)

	// 4. Reservar o ponto via HTTP
	urlReserva := fmt.Sprintf("%s/reserve-points/%s", servidor, pontoEscolhido.ID)
	req, _ := http.NewRequest("POST", urlReserva, nil)
	respReserva, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer respReserva.Body.Close()
	var resultado map[string]interface{}
	json.NewDecoder(respReserva.Body).Decode(&resultado)
	fmt.Printf("Resultado da reserva: %v\n", resultado)

	// 5. Simular viagem
	duracao := rand.Intn(10) + 10 // 3 a 7 segundos
	fmt.Printf("Simulando viagem por %d segundos...\n", duracao)
	time.Sleep(time.Duration(duracao) * time.Second)

	// 6. Liberar o ponto (cancelar reserva)
	urlCancel := fmt.Sprintf("%s/cancel-reservation", servidor)
	body := fmt.Sprintf(`{"ids":["%s"]}`, pontoEscolhido.ID)
	respCancel, err := http.Post(urlCancel, "application/json", strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer respCancel.Body.Close()
	var cancelResult map[string]interface{}
	json.NewDecoder(respCancel.Body).Decode(&cancelResult)
	fmt.Printf("Ponto liberado: %v\n", cancelResult)
}
