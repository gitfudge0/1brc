package solver

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type station struct {
	min   int
	max   int
	avg   int
	count int
}

type stationInterface interface {
	ToString() string
}

func (s *station) ToString() string {
	return fmt.Sprintf("%d | %d | %d\n", s.min, s.max, s.avg)
}

type StationList = map[string]*station

type Solver struct {
	list StationList
	wg   *sync.WaitGroup
	mu   sync.RWMutex
}
type SolverInterface interface {
	Read(key string) int
	Write(key string, value int)
}

func (s Solver) Read(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.list[key]
}

func (s Solver) Write(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list[key] = value
}

func parseAndSolve(solver *Solver, batchData []string) {
	defer solver.wg.Done()

	fmt.Println("Solving batch")
	for _, line := range batchData {
		lineData := strings.Split(line, ",")
		city := lineData[0]
		temp, _ := strconv.Atoi(lineData[1])

		solver.mu.Lock()
		if _, ok := solver.list[city]; ok {
			item := solver.list[city]

			if temp < item.min {
				item.min = temp
			}

			if temp > item.max {
				item.max = temp
			}

			item.avg = ((item.avg*item.count)+temp)/item.count + 1

			item.count += 1
		} else {
			solver.list[city] = &station{
				min:   temp,
				max:   temp,
				avg:   temp,
				count: 1,
			}
		}

		solver.mu.Unlock()
		fmt.Println(solver.list[city].ToString())
	}
}

func Solve() {
	fmt.Println("Solving")

	file, fileErr := os.OpenFile("1brc.csv", os.O_RDONLY, 777)
	if fileErr != nil {
		fmt.Println(fmt.Errorf("error opening file: %m", fileErr))
	}

	stationList := make(StationList)
	guard := make(chan struct{}, 1000)
	waitGroup := &sync.WaitGroup{}

	const batchSize = 1000
	batchData := make([]string, 0, batchSize)

	solver := Solver{
		list: stationList,
		wg:   waitGroup,
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		batchData = append(batchData, scanner.Text())

		if len(batchData) >= batchSize {
			guard <- struct{}{}
			solver.wg.Add(1)

			go func(dataCopy []string) {
				defer func() { <-guard }()
				parseAndSolve(&solver, dataCopy)
			}(batchData[:])

			batchData = make([]string, 0, batchSize)
		}
	}

	solver.wg.Wait()
	close(guard)

	for key, value := range solver.list {
		fmt.Printf("%s: %d | %d | %d\n", key, value.min, value.max, value.avg)
	}
}
