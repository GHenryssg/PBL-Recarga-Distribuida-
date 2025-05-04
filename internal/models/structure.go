package models

type PontoRecarga struct {
	ID          string `json:"id"`
	Localizacao string `json:"localizacao"`
	Disponivel  bool   `json:"disponivel"`
}

type Empresa struct {
	ID     string         `json:"id"`
	Nome   string         `json:"nome"`
	Pontos []PontoRecarga `json:"pontos"`
}

type RequisicaoVeiculo struct {
	VeiculoID string `json:"veiculo_id"`
	Bateria   int    `json:"bateria"`
	Local     string `json:"local"`
	Destino   string `json:"destino"`
}

type RespostaServidor struct {
	VeiculoID         string         `json:"veiculo_id"`
	PontosDisponiveis []PontoRecarga `json:"pontos_disponiveis"`
}

type Rota struct {
    ID     string         `json:"id"`
    Nome   string         `json:"nome"`
    Pontos []PontoRecarga `json:"pontos"`
}