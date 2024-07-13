package helpers

import (
	"etl-with-golang/models"
	"sync"
	"etl-with-golang/infra/database"
	"fmt"
	"github.com/spf13/viper"
)

var createWorker = func(wg *sync.WaitGroup, batchChan <-chan []models.Register) {
	go Worker(wg, batchChan)
}

func Worker(wg *sync.WaitGroup, batchChan <-chan []models.Register) {
	defer wg.Done()
	for batch := range batchChan {
		if err := database.DB.Create(&batch).Error; err != nil {
			fmt.Println("Error inserting batch:", err)
		}
	}
}

func SetupWorkers() (chan []models.Register, *sync.WaitGroup) {
	batchChan := make(chan []models.Register)
	var wg sync.WaitGroup

	for i := 0; i < viper.GetInt("NUM_WORKERS"); i++ {
		wg.Add(1)
		createWorker(&wg, batchChan)
	}
	return batchChan, &wg
}