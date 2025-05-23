<img width=100% src="https://capsule-render.vercel.app/api?type=waving&color=00FFFF&height=120&section=header"/>

<div align="center">  

# Projeto de Concorrência e Conectividade com Docker e Go

## Descrição do Projeto

Este projeto visa resolver o problema de "ansiedade de autonomia" em veículos elétricos (VEs) no Brasil, desenvolvendo um sistema distribuído para planejamento e reserva antecipada de múltiplos pontos de recarga ao longo de rotas entre cidades e estados. O sistema permite que usuários realizem reservas atômicas em pontos de diferentes empresas através de uma API padronizada.

## Fundionalidades Principais

- **Para usuários (veículos)**:

Consulta de rotas disponíveis com pontos de recarga
Reserva atômica de múltiplos pontos em uma rota
Simulação de viagem com monitoramento de bateria
Cancelamento de reservas

- **Para servidores (empresas)**:

Gerenciamento de pontos de recarga locais
Comunicação com outros servidores via REST e MQTT
Coordenação distribuída para reservas atômicas
Rollback automático em caso de falha

## Tecnologias Utilizadas

- **Linguagem**: Go (Golang)
- **Contêinerização**: Docker e Docker Compose
- **Comunicação entre servidores**: API REST
- **Comunicação com veículos**: Protocolo MQTT (IoT)
- **Framework web**: Gin
- **Broker MQTT**: Eclipse Paho

## Como Executar o Projeto

### Pré-requisitos

- Docker instalado
- Docker Compose instalado

### Passos para Execução

1. Clone o repositório do projeto
2. Navegue até o diretório do projeto
3. Execute o seguinte comando:

```bash
docker-compose up --build
```

Isso irá construir e iniciar todos os serviços definidos no docker-compose.yml.

## Estrutura de Arquivos

- `main.go`: Ponto de entrada do servidor

- `client`: Implementação do cliente MQTT

- `controllers`: Handlers HTTP para endpoints da API

- `services`: Lógica de negócio para pontos e rotas

- `models`: Estruturas de dados do sistema

- `config`: Configurações e variáveis de ambiente

- `docker`: Configurações Docker


## Protocolos de Comunicação
**Entre servidores (REST)**:
GET /routes - Lista todas as rotas disponíveis

POST /reserve-points/{ids} - Reserva múltiplos pontos

POST /cancel-reservation/{ids} - Cancela reservas

**Com veículos (MQTT)**:
Tópico veiculos/solicitacao - Solicitações de pontos disponíveis

Tópico veiculos/resposta/{id} - Respostas para veículos específicos

Tópico empresa/{nome}/reserva - Comunicação entre empresas

## Tratamento de Concorrência e Consistência
O sistema utiliza várias técnicas para garantir operações confiáveis:

Reservas atômicas: Verifica disponibilidade de todos os pontos antes de confirmar reservas

Rollback distribuído: Cancela todas as reservas se qualquer ponto ficar indisponível

Mutexes: Protegem estruturas de dados compartilhadas nos servidores

Canais: Gerenciam comunicação assíncrona no cliente MQTT

Timeouts: Garantem que operações remotas não fiquem bloqueadas indefinidamente

## Exemplo de Fluxo
1. Veículo consulta rotas disponíveis via HTTP GET /routes

2. Seleciona pontos de recarga ao longo da rota

3. Envia requisição de reserva atômica via HTTP POST /reserve-points

4. Servidor coordena com outros servidores via REST/MQTT

5. Verifica disponibilidade em todos os pontos

6. Reserva todos ou nenhum (atomicidade)

7. Se bem-sucedido, veículo inicia viagem simulada

8. Ao final, libera pontos via POST /cancel-reservation


## Monitoramento e Logs

Monitoramento e Logs

O sistema gera logs detalhados para cada componente:
- Servidores registram todas as operações de reserva, cancelamento e comunicação entre empresas
- Pontos de recarga registram status de disponibilidade e histórico de uso
- Veículos registram consumo de bateria e operações de recarga
- Broker MQTT registra todas as mensagens trocadas entre os componentes
- Os dados são armazenados em:
- Arquivos JSON no volume compartilhado dados_recarga
- Banco de dados temporal para métricas de desempenho
- Stream de logs centralizado para monitoramento em tempo real

## Personalização

O sistema permite customização através de:
- **docker-compose.yml - Para ajustar**:

Número de servidores de empresas
Quantidade de veículos simulados
Configurações de rede entre containers

- **Arquivos de configuração - Para modificar**:

config/config.go - Parâmetros de tempo de recarga
models/carro.go - Consumo de bateria e comportamento dos veículos
.env - URLs de conexão e configurações por ambiente

- **Variáveis de ambiente - Para controlar**:

Níveis de log (DEBUG, INFO, WARN)
Timeouts de comunicação
Políticas de retentativa
Configurações específicas por empresa



## Relatório


## **Alunos(as)**

<table align='center'>
<tr> 
  <td align="center">
    <a href="https://github.com/LuisMarioRC">
      <img src="https://avatars.githubusercontent.com/u/142133059?v=4" width="100px;" alt=""/>
    </a>
    <br /><sub><b><a href="https://github.com/LuisMarioRC">Luis Mario</a></b></sub><br />👨💻
  </td>
  <td align="center">
    <a href="https://github.com/laizagordiano">
      <img src="https://avatars.githubusercontent.com/u/132793645?v=4" width="100px;" alt=""/>
    </a>
    <br /><sub><b><a href="https://github.com/laizagordiano">Laiza Gordiano</a></b></sub><br />👨💻
  </td>
  <td align="center">
    <a href="https://github.com/GHenryssg">
      <img src=https://avatars.githubusercontent.com/u/142272107?v=4" width="100px;" alt=""/>
    </a>
    <br /><sub><b><a href="https://github.com/GHenryssg">Gabriel Henry</a></b></sub><br />👨💻
  </td>
</tr>

</table>


<img width=100% src="https://capsule-render.vercel.app/api?type=waving&color=00FFFF&height=120&section=footer"/>

<div align="center"> 