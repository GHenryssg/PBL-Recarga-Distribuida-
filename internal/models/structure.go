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

type Rota struct {
    ID     string         `json:"id"`
    Nome   string         `json:"nome"`
    Pontos []PontoRecarga `json:"pontos"`
}

type PlanejamentoViagem struct {
	Origem  string `json:"origem"`
	Destino string `json:"destino"`
}

type PontoReserva struct {
	ID      string `json:"id"`
	Empresa string `json:"empresa"` // ainda sem uso
}

type SequenciaReserva struct {
	Pontos []PontoReserva `json:"pontos"`
}
