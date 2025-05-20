package config

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

var (
	NomeEmpresa       string
	Porta             string
	MQTTBrokerURL     string
	URLsEmpresas      map[string]string
	EmpresaNomeParaID map[string]string
)

func CarregarVariaveis() {
	// Tente carregar o .env no diretório atual, mas NÃO sobrescreva variáveis já definidas no ambiente
	err := godotenv.Load(".env")
	if err != nil {
		// Se não encontrar, tente carregar o .env no caminho do Docker
		err = godotenv.Load("/app/.env")
		if err != nil {
			log.Println("⚠️  Arquivo .env não encontrado, usando variáveis do sistema")
		}
	}

	NomeEmpresa = os.Getenv("NOME_EMPRESA")
	Porta = os.Getenv("PORTA")
	MQTTBrokerURL = os.Getenv("MQTT_BROKER")

	if _, err := net.LookupHost("mqtt_broker"); err != nil {
		log.Println("⚠️  Host mqtt_broker não resolvível, usando localhost")
		MQTTBrokerURL = "tcp://localhost:1883"
	}

	URLsEmpresas = map[string]string{
		"empresa_a": os.Getenv("EMPRESA_A_URL"),
		"empresa_b": os.Getenv("EMPRESA_B_URL"),
		"empresa_c": os.Getenv("EMPRESA_C_URL"),
	}

	EmpresaNomeParaID = map[string]string{
		"empresa_a": "1",
		"empresa_b": "2",
		"empresa_c": "3",
	}
}
