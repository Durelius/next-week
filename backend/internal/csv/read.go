package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func Read(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 2. Create a CSV reader
	reader := csv.NewReader(file)

	// 3. Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
