package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"etl-with-golang/repository"
	"etl-with-golang/helpers"
	"fmt"
)

// GetImportReport retrieves a report for a specific import operation.
//
// GET /api/v1/import-report
//
// Query Parameters:
//   - importacaoId: UUID of the import operation (required)
//
// Responses:
//   - 200 OK: Report retrieved successfully
//     {
//       "totalRows": <int>,
//       "invalidCPFCount": <int>,
//       "invalidLojaMaisFrequenteCNPJCount": <int>,
//       "invalidLojaUltimaCompraCNPJCount": <int>
//     }
//   - 400 Bad Request: Invalid importacaoId
//   - 404 Not Found: Import operation records not found
//   - 500 Internal Server Error: Server-side processing error
//
// The report includes the total number of rows processed, and counts of invalid
// CPFs, invalid Loja Mais Frequente CNPJs, and invalid Loja Ultima Compra CNPJs.
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

	if totalRows == 0 {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "No records found for this importacaoId"})
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