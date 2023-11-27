package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	// open csv
	file, err := os.Open("input.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// pre process
	fmt.Println(data)

	// process data
	for i, row := range data {
		if i > 0 {
			for j, col := range row {
				wg.Add(1)
				go func(i, j int, data *string, col string) {
					defer wg.Done()
					mu.Lock()
					if j == 0 {
						*data = strings.ToUpper(col)
					} else if j == 2 {
						*data = "Mr." + col
					}
					mu.Unlock()
				}(i, j, &data[i][j], col)
			}
		}
	}
	wg.Wait()

	// post process
	fmt.Println(data)

	// Write the CSV data
	file2, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer file2.Close()

	writer := csv.NewWriter(file2)
	for _, row := range data {
		writer.Write(row)
	}

}
