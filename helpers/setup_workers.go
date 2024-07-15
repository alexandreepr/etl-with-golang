package helpers

import (
	"etl-with-golang/models"
	"sync"
	"etl-with-golang/infra/database"
	"etl-with-golang/infra/logger"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

func ProcessLineWorker(id int, tasks <-chan string, results chan<- models.Register, importId uuid.UUID, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range tasks {
		record, err := ProcessLine(line, importId)
		if err != nil {
			logger.Errorf("Worker %d: error processing line: %v\n", id, err)
			continue
		}

		results <- record
	}
}

func BatchWorker(wg *sync.WaitGroup, results <-chan models.Register, batchSize int) {
	defer wg.Done()
	var records []models.Register

	for record := range results {
		records = append(records, record)
		if len(records) >= batchSize {
			err := database.DB.Transaction(func(tx *gorm.DB) error {
				return tx.CreateInBatches(records, batchSize).Error
			})
			if err != nil {
				logger.Errorf("Error inserting batch: %v", err)
			}
			records = nil
		}
	}

	if len(records) > 0 {
		err := database.DB.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(records, batchSize).Error
		})
		if err != nil {
			logger.Errorf("Error inserting batch: %v", err)
		}
	}
}