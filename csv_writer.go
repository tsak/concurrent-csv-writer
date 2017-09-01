// Package ccsv provides a "thread" safe way of writing to CSV files
package ccsv

import (
	"encoding/csv"
	"os"
	"sync"
)

// CsvWriter holds pointers to a Mutex, csv.Writer and the underlying CSV file
type CsvWriter struct {
	mutex     *sync.Mutex
	csvWriter *csv.Writer
	file      *os.File
}

// NewCsvWriter creates a CSV file and returns a CsvWriter
func NewCsvWriter(fileName string) (*CsvWriter, error) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(csvFile)
	return &CsvWriter{csvWriter: w, mutex: &sync.Mutex{}, file: csvFile}, nil
}

// Write a single row to a CSV file
func (w *CsvWriter) Write(row []string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.csvWriter.Write(row)
}

// WriteAll writes multiple rows to a CSV file
func (w *CsvWriter) WriteAll(records [][]string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.csvWriter.WriteAll(records)
}

// Flush writes any pending rows
func (w *CsvWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.csvWriter.Flush()
	return w.csvWriter.Error()
}

// Close CSV file for writing
// Implicitly calls Flush() before
func (w *CsvWriter) Close() error {
	err := w.Flush()
	if err != nil {
		return err
	}
	return w.file.Close()
}
