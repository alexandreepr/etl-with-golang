package helpers

import (
	"testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProcessLine(t *testing.T) {
	importId := uuid.New()
	line := "288.663.049-72 1 0 2021-12-01 100.00 150.00 33.014.556/0098-19 33.014.556/0098-19"

	record, err := ProcessLine(line, importId)
	assert.NoError(t, err)
	assert.Equal(t, importId, record.ImportacaoId)
	assert.Equal(t, "28866304972", record.CPF)
	assert.Equal(t, true, record.CPFValido)
	assert.Equal(t, true, record.Private)
	assert.Equal(t, false, record.Incompleto)
	assert.Equal(t, "2021-12-01", record.DataUltimaCompra)
	assert.Equal(t, "100.00", record.TicketMedio)
	assert.Equal(t, "150.00", record.TicketUltimaCompra)
	assert.Equal(t, "33014556009819", record.LojaMaisFrequente)
	assert.Equal(t, true, record.LojaMaisFrequenteCNPJValido)
	assert.Equal(t, "33014556009819", record.LojaUltimaCompra)
	assert.Equal(t, true, record.LojaUltimaCompraCNPJValido)
}

func TestSanitizeString(t *testing.T) {
	assert.Equal(t, "12345678901", SanitizeString("123.456.789-01"))
	assert.Equal(t, "NULL", SanitizeString("NULL"))
	assert.Equal(t, "12345678901", SanitizeString("12345678901"))
}

func TestParseToBool(t *testing.T) {
	assert.Equal(t, true, ParseToBool("1"))
	assert.Equal(t, false, ParseToBool("0"))
	assert.Equal(t, false, ParseToBool("invalid"))
}
