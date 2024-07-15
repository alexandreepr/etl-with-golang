package repository

import (
	"etl-with-golang/infra/database"
	"etl-with-golang/infra/logger"
	"etl-with-golang/models"
	"github.com/google/uuid"
)

func CountCPFValidoFalse(importacaoId uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Register{}).Where("importacao_id = ? AND cpf_valido = ?", importacaoId, false).Count(&count).Error
	if err != nil {
		logger.Errorf("error counting CPFValido false: %v", err)
	}
	return count, err
}

func CountLojaMaisFrequenteCNPJValidoFalse(importacaoId uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Register{}).Where("importacao_id = ? AND loja_mais_frequente_cnpj_valido = ?", importacaoId, false).Count(&count).Error
	if err != nil {
		logger.Errorf("error counting LojaMaisFrequenteCNPJValido false: %v", err)
	}
	return count, err
}

func CountLojaUltimaCompraCNPJValidoFalse(importacaoId uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Register{}).Where("importacao_id = ? AND loja_ultima_compra_cnpj_valido = ?", importacaoId, false).Count(&count).Error
	if err != nil {
		logger.Errorf("error counting LojaUltimaCompraCNPJValido false: %v", err)
	}
	return count, err
}

func CountImportTotalRows(importacaoId uuid.UUID) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Register{}).Where("importacao_id = ?", importacaoId).Count(&count).Error
	return count, err
}