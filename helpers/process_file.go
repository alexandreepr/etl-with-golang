package helpers

import (
	"bufio"
	"fmt"
	"strings"
	"etl-with-golang/models"
	"mime/multipart"
	"strconv"
	"regexp"
	"github.com/paemuri/brdoc"
	"github.com/spf13/viper"
)

func ProcessFile(file multipart.File, batchChan chan<- []models.Register) error {
	scanner := bufio.NewScanner(file)
	var records []models.Register
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Skip header line
		if lineNumber == 1 {
			continue
		}

		record, err := ProcessLine(line, lineNumber)
		if err != nil {
			return err
		}

		records = append(records, record)

		if len(records) == viper.GetInt("BATCH_SIZE") {
			batchChan <- records
			records = nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(records) > 0 {
		batchChan <- records
	}

	return nil
}

func ProcessLine(line string, lineNumber int) (models.Register, error) {
	fields := strings.Fields(line)
	
	if len(fields) != 8 {
		return models.Register{}, fmt.Errorf("incorrect file format on line %d", lineNumber)
	}

	// TODO: add ImportId
	return models.Register{
		CPF:                		  SanitizeString(fields[0]),
		CPFValido:          		  brdoc.IsCPF(fields[0]),
		Private:            		  ParseToBool(fields[1]),
		Incompleto:         		  ParseToBool(fields[2]),
		DataUltimaCompra:   		  fields[3],
		TicketMedio:        		  fields[4],
		TicketUltimaCompra: 		  fields[5],
		LojaMaisFrequente:  		  SanitizeString(fields[6]),
		LojaMaisFrequenteCNPJValido:  brdoc.IsCNPJ(fields[6]),
		LojaUltimaCompra:   		  SanitizeString(fields[7]),
		LojaUltimaCompraCNPJValido:   brdoc.IsCNPJ(fields[7]),
	}, nil
}

func SanitizeString(dirtyString string) string {
	if dirtyString == "NULL" {
		return dirtyString
	}
	
	// Removes all the non digit characters from the string
	reg := regexp.MustCompile(`[^\d]`)
	return reg.ReplaceAllString(dirtyString, "")
}

func ParseToBool(str string) bool {
	value, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return value == 1
}