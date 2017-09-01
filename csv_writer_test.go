package ccsv

import (
	"bufio"
	"os"
	"strconv"
	"testing"
)

func getCsvWriter() *CsvWriter {
	csvWriter, err := NewCsvWriter("test.csv")
	if err != nil {
		panic("could not open test.csv for writing")
	}
	return csvWriter
}

func rowCount() int {
	file, _ := os.Open("test.csv")
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	rowCount := 0
	for fileScanner.Scan() {
		rowCount++
	}
	return rowCount
}

func deleteCsvFile() {
	os.Remove("test.csv")
}

var (
	row = []string{"A", "row", "in", "a", "CSV", "file"}

	rows = [][]string{
		{"First", "of", "three", "rows"},
		{"Second", "of", "three", "rows"},
		{"Third", "of", "three", "rows"},
	}
)

func TestNewCsvWriter(t *testing.T) {
	csvWriter, err := NewCsvWriter("test.csv")
	if err != nil {
		t.Error("should open a CSV file for writing")
	}

	err = csvWriter.Close()
	if err != nil {
		t.Error("should have closed CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_Write(t *testing.T) {
	csvWriter := getCsvWriter()

	t.Run("multiple single row writes", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run("single row write "+strconv.Itoa(i), func(t *testing.T) {
				err := csvWriter.Write(row)
				if err != nil {
					t.Error("should have written a row to the CSV file")
				}
			})
		}
	})

	csvWriter.Flush()
	csvWriter.Close()

	if ok := rowCount() != 10; ok {
		t.Error("expected to see 10 rows in CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_WriteAll(t *testing.T) {
	csvWriter := getCsvWriter()
	defer csvWriter.Close()

	t.Run("multiple multi row writes", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run("multi row write "+strconv.Itoa(i), func(t *testing.T) {
				err := csvWriter.WriteAll(rows)
				if err != nil {
					t.Error("should have written a row to the CSV file")
				}
			})
		}
	})

	csvWriter.Flush()
	csvWriter.Close()

	if ok := rowCount() != 30; ok {
		t.Error("expected to see 30 rows in CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_Flush(t *testing.T) {
	csvWriter := getCsvWriter()

	csvWriter.Write(row)
	err := csvWriter.Flush()
	if err != nil {
		t.Error("should have flushed output buffers for CSV file")
	}
	csvWriter.Close()

	if ok := rowCount() != 1; ok {
		t.Error("expected to see one row in CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_Close(t *testing.T) {
	csvWriter := getCsvWriter()

	err := csvWriter.Close()
	if err != nil {
		t.Error("expected to close underlying file io stream")
	}

	deleteCsvFile()
}
