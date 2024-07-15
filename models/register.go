package models

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type Register struct {
	gorm.Model
	ImportacaoId       			 uuid.UUID `gorm:"column:importacao_id; type:uuid"`
	CPF                			 string `gorm:"column:cpf"`
	CPFValido          			 bool   `gorm:"column:cpf_valido"`
	Private            			 bool   `gorm:"column:private"`
	Incompleto         			 bool   `gorm:"column:incompleto"`
	DataUltimaCompra   			 string `gorm:"column:data_ultima_compra"`
	TicketMedio        			 string `gorm:"column:ticket_medio"`
	TicketUltimaCompra 			 string `gorm:"column:ticket_ultima_compra"`
	LojaMaisFrequente  			 string `gorm:"column:loja_mais_frequente"`
	LojaMaisFrequenteCNPJValido  bool   `gorm:"column:loja_mais_frequente_cnpj_valido"`
	LojaUltimaCompra   			 string `gorm:"column:loja_ultima_compra"`
	LojaUltimaCompraCNPJValido   bool   `gorm:"column:loja_ultima_compra_cnpj_valido"`
}
