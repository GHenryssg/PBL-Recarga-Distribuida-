package database

import (
	"github.com/GHenryssg/PBL-Recarga-Distribuida-/internal/models"
)

var Empresas = []models.Empresa{
    {
        ID:   "1",
        Nome: "HenryVolt",
        Pontos: []models.PontoRecarga{
            {ID: "1", Localizacao: "Feira de Santana", Disponivel: true},
            {ID: "2", Localizacao: "Salvador", Disponivel: true},
        },
    },
    {
        ID:   "2",
        Nome: "LaizaCharge",
        Pontos: []models.PontoRecarga{
            {ID: "3", Localizacao: "Valente", Disponivel: true},
            {ID: "4", Localizacao: "Santaluz", Disponivel: true},
        },
    },
    {
        ID:   "3",
        Nome: "MarioPower",
        Pontos: []models.PontoRecarga{
            {ID: "5", Localizacao: "Amargosa", Disponivel: true},
        },
    },
}

var Pontos = []models.PontoRecarga{
	{ID: "1", Localizacao: "Feira de Santana", Disponivel: true},
	{ID: "2", Localizacao: "Salvador", Disponivel: true},
	{ID: "3", Localizacao: "Valente", Disponivel: true},
	{ID: "4", Localizacao: "Santaluz", Disponivel: true},
	{ID: "5", Localizacao: "Amargosa", Disponivel: true},
}

var Rotas = []models.Rota{
	{
		ID:   "1",
		Nome: "Rota Feira de Santana - Salvador",
		Pontos: []models.PontoRecarga{
			{ID: "1", Localizacao: "Feira de Santana", Disponivel: true},
			{ID: "2", Localizacao: "Salvador", Disponivel: true},
		},
	},
	{
		ID:   "2",
		Nome: "Rota Salvador - Camaçari",
		Pontos: []models.PontoRecarga{
			{ID: "2", Localizacao: "Salvador", Disponivel: true},
			{ID: "3", Localizacao: "Camaçari", Disponivel: true},
		},
	},
	{
		ID:   "3",
		Nome: "Rota Feira de Santana  - Santaluz",
		Pontos: []models.PontoRecarga{
			{ID: "1", Localizacao: "Feira de Santana", Disponivel: true},
			{ID: "3", Localizacao: "Valente", Disponivel: true},
			{ID: "4", Localizacao: "Santaluz", Disponivel: true},
		},
	},
	{
		ID:   "4",
		Nome: "Rota Salvador - Feira de Santana",
		Pontos: []models.PontoRecarga{
			{ID: "2", Localizacao: "Salvador", Disponivel: true},
			{ID: "5", Localizacao: "Amargosa", Disponivel: true},
			{ID: "1", Localizacao: "Feira de Santana", Disponivel: true},
		},
	},
	{
		ID:   "5",
		Nome: "Rota Valente - Amargosa",
		Pontos: []models.PontoRecarga{
			{ID: "3", Localizacao: "Valente", Disponivel: true},
			{ID: "4", Localizacao: "Santaluz", Disponivel: true},
			{ID: "5", Localizacao: "Amargosa", Disponivel: true},
		},
	},
}