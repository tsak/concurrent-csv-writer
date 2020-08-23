package ccsv

import (
	"bufio"
	"os"
	"sync"
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

	err = csvWriter.Close()
	if err == nil {
		t.Error("should have complained about already closed CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_Write(t *testing.T) {
	csvWriter := getCsvWriter()

	wg := sync.WaitGroup{}
	t.Run("multiple single row writes", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				err := csvWriter.Write(row)
				if err != nil {
					t.Error("should have written a row to the CSV file")
				}
			}(i)
		}
	})

	wg.Wait()

	if err := csvWriter.Close(); err != nil {
		t.Errorf("expected no error when closing CSV file, got %s", err)
	}

	if ok := rowCount() != 1000; ok {
		t.Error("expected to see 1000 rows in CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_WriteAll(t *testing.T) {
	csvWriter := getCsvWriter()

	wg := sync.WaitGroup{}
	t.Run("multiple multi row writes", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := csvWriter.WriteAll(rows)
				if err != nil {
					t.Error("should have written a row to the CSV file")
				}
			}()
		}
	})

	wg.Wait()

	if err := csvWriter.Close(); err != nil {
		t.Errorf("expected no error when closing CSV file, got %s", err)
	}

	if ok := rowCount() != 3000; ok {
		t.Error("expected to see 3000 rows in CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_Flush(t *testing.T) {
	csvWriter := getCsvWriter()

	if err := csvWriter.WriteAll(rows); err != nil {
		t.Errorf("expected no error when writing csv file, got: %s", err)
	}

	if err := csvWriter.Flush(); err != nil {
		t.Error("should have flushed output buffers for CSV file")
	}

	if err := csvWriter.Close(); err != nil {
		t.Errorf("expected no error when closing CSV file, got %s", err)
	}

	if ok := rowCount() != 3; ok {
		t.Error("expected to see three rows in CSV file")
	}

	deleteCsvFile()
}

func TestCsvWriter_Close(t *testing.T) {
	csvWriter := getCsvWriter()

	if err := csvWriter.Close(); err != nil {
		t.Errorf("expected to close underlying file io stream without error, got: %s", err)
	}

	deleteCsvFile()
}

func BenchmarkCsvWriter_Write(b *testing.B) {
	csvWriter := getCsvWriter()

	b.Run("nil-slice", func(b *testing.B) {
		wg := sync.WaitGroup{}
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
				go func() {
					defer wg.Done()
					csvWriter.Write(nil)
				}()
			}
		})
		wg.Wait()
	})

	b.Run("row", func(b *testing.B) {
		wg := sync.WaitGroup{}
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
				go func() {
					defer wg.Done()
					csvWriter.Write(row)
				}()
			}
		})
		wg.Wait()
	})

	csvWriter.Close()
	deleteCsvFile()
}
