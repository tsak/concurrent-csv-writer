// Package ccsv provides a concurrency safe way of writing to CSV files from multiple go routines
package ccsv

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
)

// CsvWriter encapsulates the underlying os.File, a csv.Writer, a data channel, a done channel and
// a flag that determines if a writer is already closed
type CsvWriter struct {
	file   *os.File
	writer *csv.Writer
	data   chan []string
	done   chan bool
	closed bool
}

var ErrWriterClosed = errors.New("ccsv: writer already closed")

// NewCsvWriter creates a CSV file and returns a pointer to a new CsvWriter
// Will return an error if the underlying os.Create fails for some reason
func NewCsvWriter(filename string) (*CsvWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	ccsvwriter := CsvWriter{
		file:   file,
		writer: csv.NewWriter(file),
		data:   make(chan []string),
		done:   make(chan bool),
	}

	go ccsvwriter.loop()

	return &ccsvwriter, nil
}

// loop handles writes via channel and closing the writer if the data channel is closed
func (c *CsvWriter) loop() {
	for {
		select {
		case msg, ok := <-c.data:
			// Data channel was closed, let's wrap things up
			if !ok {
				c.closed = true
				// Flush any pending writes
				c.writer.Flush()
				if err := c.writer.Error(); err != nil {
					log.Println("Error", err)
				}
				if err := c.file.Close(); err != nil {
					log.Println(err)
				}
				c.done <- true
				return
			}
			if err := c.writer.Write(msg); err != nil {
				log.Println(err)
			}
			if err := c.writer.Error(); err != nil {
				log.Println("Error", err)
			}
		}
	}
}

// Write a single row to a CSV file
// Will return a ErrWriterClosed error if the CsvWriter was closed before
func (c *CsvWriter) Write(row []string) (err error) {
	defer func() {
		// Any pending writes once the data channel is closed can result in a panic, this will recover from it
		if recover() != nil {
			err = ErrWriterClosed
		}
	}()

	if c.closed {
		return ErrWriterClosed
	}
	c.data <- row
	return
}

// WriteAll writes multiple rows to a CSV file
// Will return a ErrWriterClosed error if the CsvWriter was closed before
func (c *CsvWriter) WriteAll(rows [][]string) error {
	for _, row := range rows {
		if err := c.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// Flush writes any pending rows
// Will return a ErrWriterClosed error if the CsvWriter was closed before
func (c *CsvWriter) Flush() error {
	if c.closed {
		return ErrWriterClosed
	}
	c.writer.Flush()
	return c.writer.Error()
}

// Close CSV file for writing
// Implicitly calls Flush() before
// Will return a ErrWriterClosed error if the CsvWriter was closed before
func (c *CsvWriter) Close() error {
	if c.closed {
		return ErrWriterClosed
	}
	close(c.data)
	<-c.done
	close(c.done)
	return nil
}
