package helpers

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"etl-with-golang/models"
	"etl-with-golang/infra/logger"
	"mime/multipart"
	"strconv"
	"regexp"
	"github.com/paemuri/brdoc"
	"github.com/spf13/viper"
	"github.com/google/uuid"
)

func ProcessFile(file multipart.File) (string, error) {
	scanner := bufio.NewScanner(file)
	tasks := make(chan string, 100)
	results := make(chan models.Register, 100)
	var wg sync.WaitGroup
	importId := uuid.New()

	numWorkers := viper.GetInt("NUM_WORKERS")
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go ProcessLineWorker(i, tasks, results, importId, &wg)
	}

	var batchWg sync.WaitGroup
	batchSize := viper.GetInt("BATCH_SIZE")
	numBatchWorkers := viper.GetInt("NUM_WORKERS")
	for i := 0; i < numBatchWorkers; i++ {
		batchWg.Add(1)
		go BatchWorker(&batchWg, results, batchSize)
	}

	go func() {
		lineNumber := 0
		for scanner.Scan() {
			line := scanner.Text()

			// Skip header line
			if lineNumber == 0 {
				lineNumber++
				continue
			}

			tasks <- line
		}
		close(tasks)

		if err := scanner.Err(); err != nil {
			logger.Errorf("Scanner error: %v\n", err)
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// Wait for all batch workers to complete
	batchWg.Wait()

	return importId.String(), nil
}

func ProcessLine(line string, importId uuid.UUID) (models.Register, error) {
	fields := strings.Fields(line)

	if len(fields) != 8 {
		return models.Register{}, fmt.Errorf("incorrect file format")
	}

	return models.Register{
		ImportacaoId:                 importId,
		CPF:                          SanitizeString(fields[0]),
		CPFValido:                    brdoc.IsCPF(fields[0]),
		Private:                      ParseToBool(fields[1]),
		Incompleto:                   ParseToBool(fields[2]),
		DataUltimaCompra:             fields[3],
		TicketMedio:                  fields[4],
		TicketUltimaCompra:           fields[5],
		LojaMaisFrequente:            SanitizeString(fields[6]),
		LojaMaisFrequenteCNPJValido:  brdoc.IsCNPJ(fields[6]),
		LojaUltimaCompra:             SanitizeString(fields[7]),
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