package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func getRandNumber(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func getRandTemperature() int {
	return getRandNumber(-99, 99)
}

func getRandCityIndex() int {
	return getRandNumber(1, 9998)
}

func getCitiesList() []string {
	file, err := os.OpenFile("cities.txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(fmt.Errorf("could not open cities file: %m", err))
	}

	cities := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		city := scanner.Text()
		cities = append(cities, city)
	}

	return cities
}

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
	if err != nil {
		log.Fatal(fmt.Errorf("could not open cities file: %m", err))
	}
	processor.mu.Unlock()
}

func createData() {
	csvFile, csvFileErr := os.OpenFile("1brc.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	if csvFileErr != nil {
		log.Fatal(fmt.Errorf("error reading file %m", csvFileErr))
	}
	defer csvFile.Close()

	cities := getCitiesList()

	waitGroup := &sync.WaitGroup{}
	writer := csv.NewWriter(bufio.NewWriter(csvFile))

	processor := Processor{
		writer: writer,
		file:   csvFile,
		wg:     waitGroup,
	}
	defer processor.writer.Flush()

	guard := make(chan struct{}, 100)
	batchSize := 10000
	batchData := make([][]string, 0, batchSize)

	var data []string

	for i := 0; i < 1_000_000_000; i++ {
		cityIndex := getRandCityIndex()
		temp := strconv.Itoa(getRandTemperature())
		data = []string{cities[cityIndex], temp}

		batchData = append(batchData, data)

		if len(batchData) >= batchSize {
			guard <- struct{}{}
			processor.wg.Add(1)

			go func(dataCopy [][]string) {
				defer func() { <-guard }()
				appendString(&processor, dataCopy)
			}(batchData[:])

			batchData = make([][]string, 0, batchSize)
		}
	}

	processor.wg.Wait()
}

func main() {
	startTime := time.Now()

	fmt.Println("Hold on while the file is being generated...")
	createData()

	fmt.Println("done")
	fmt.Printf("Time taken: %s\n", time.Since(startTime))
}
