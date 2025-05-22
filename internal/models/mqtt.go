package models

type RequisicaoVeiculo struct {
	VeiculoID string `json:"veiculo_id"`
	Bateria   int    `json:"bateria"`
	Local     string `json:"local"`
	Destino   string `json:"destino"`
}

// Modelos de mensagem usados na comunicação MQTT
// RequisicaoVeiculo: mensagem enviada pelo cliente para solicitar reserva/cancelamento
// RespostaServidor: resposta do backend para o cliente, com pontos disponíveis ou status da reserva
// Essas estruturas são serializadas em JSON e trafegam nos tópicos MQTT

type RespostaServidor struct {
	VeiculoID         string         `json:"veiculo_id"`
	PontosDisponiveis []PontoRecarga `json:"pontos_disponiveis"`
}
