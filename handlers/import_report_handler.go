package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"etl-with-golang/repository"
	"etl-with-golang/helpers"
	"fmt"
)

func GetImportReport(ctx *gin.Context) {
	importacaoIdStr := ctx.Query("importacaoId")
	fmt.Print(importacaoIdStr)
	importacaoId, err := uuid.Parse(importacaoIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid importacaoId"})
		return
	}

	totalRows, err := repository.CountImportTotalRows(importacaoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total rows"})
		return
	}

	invalidCPFCount, err := repository.CountCPFValidoFalse(importacaoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count invalid CPFs"})
		return
	}

	invalidLojaMaisFrequenteCNPJCount, err := repository.CountLojaMaisFrequenteCNPJValidoFalse(importacaoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count invalid Loja Mais Frequente CNPJ"})
		return
	}

	invalidLojaUltimaCompraCNPJCount, err := repository.CountLojaUltimaCompraCNPJValidoFalse(importacaoId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count invalid Loja Ultima Compra CNPJ"})
		return
	}

	report := helpers.ReportResponse{
		TotalRows:                         totalRows,
		InvalidCPFCount:                   invalidCPFCount,
		InvalidLojaMaisFrequenteCNPJCount: invalidLojaMaisFrequenteCNPJCount,
		InvalidLojaUltimaCompraCNPJCount:  invalidLojaUltimaCompraCNPJCount,
	}

	ctx.JSON(http.StatusOK, report)
}