package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

type Processor struct {
	writer *csv.Writer
	file   *os.File
	wg     *sync.WaitGroup
	mu     sync.Mutex
}

func appendString(processor *Processor, data [][]string) {
	defer processor.wg.Done()

	processor.mu.Lock()
	err := processor.writer.WriteAll(data)
	processor.mu.Unlock()
	if err != nil {
		fmt.Println(fmt.Errorf("error %m", err))
	}
	fmt.Print(".")
}

func createData() {
	file, err := os.OpenFile("1brc.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		errString := fmt.Errorf("error reading file %m", err)
		fmt.Println(errString)
		os.Exit(2)
	}
	defer file.Close()

	data := [][]string{
		{"Hamburg", "12.0"},
		{"Bulawayo", "8.9"},
		{"Palembang", "38.8"},
		{"St. John's", "15.2"},
		{"Cracow", "12.6"},
		{"Bridgetown", "26.9"},
		{"Istanbul", "6.2"},
		{"Roseau", "34.4"},
		{"Conakry", "31.2"},
		{"Istanbul", "23.0"},
	}

	waitGroup := &sync.WaitGroup{}
	writer := csv.NewWriter(file)

	processor := Processor{
		writer: writer,
		file:   file,
		wg:     waitGroup,
	}

	guard := make(chan struct{}, 100)
	batchSize := 10000
	batchData := make([][]string, 0, batchSize)

	for i := 0; i < 100000000; i++ {
		batchData = append(batchData, data...)

		if len(batchData) >= batchSize {
			fmt.Printf("Batch size met at %d\n", i)
			guard <- struct{}{}
			processor.wg.Add(1)

			go func() {
				defer func() { <-guard }()
				appendString(&processor, batchData)
			}()

			batchData = make([][]string, 0, batchSize)
		}
	}

	processor.wg.Wait()
}

func main() {
	startTime := time.Now()

	createData()

	fmt.Println("done")
	fmt.Printf("Time taken: %s\n", time.Since(startTime))
}
