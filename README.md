# Concurrent CSV writer

[![Go Report Card](https://goreportcard.com/badge/github.com/tsak/concurrent-csv-writer)](https://goreportcard.com/report/github.com/tsak/concurrent-csv-writer)
[![GoDoc](https://godoc.org/github.com/tsak/concurrent-csv-writer?status.svg)](https://godoc.org/github.com/tsak/concurrent-csv-writer)

A thread-safe way of concurrent writes to a CSV file in Go. Order of rows is **NOT** guaranteed.

Inspired by a [blog post](http://www.markhneedham.com/blog/2017/01/31/go-multi-threaded-writing-csv-file/) by [Mark Needham](http://www.markhneedham.com).

## Usage

```go
import (
    "github.com/tsak/concurrent-csv-writer"
)
```
    
```bash
go get "github.com/tsak/concurrent-csv-writer"
```    
    
## Example

```go
    package main
    
    import (
        "github.com/tsak/concurrent-csv-writer"
        "strconv"
    )
    
    func main() {
        // Create `sample.csv` in current directory
        csv, err := ccsv.NewCsvWriter("sample.csv")
        if err != nil {
            panic("Could not open `sample.csv` for writing")
        }
    
        // Flush pending writes and close file upon exit of main()
        defer csv.Close()
    
        count := 99
    
        done := make(chan bool)
    
        for i := count; i > 0; i-- {
            go func(i int) {
                csv.Write([]string{strconv.Itoa(i), "bottles", "of", "beer"})
                done <- true
            }(i)
        }
    
        for i := 0; i < count; i++ {
            <-done
        }
    }
```

### Output

Notice the lack of order of entries.

```csv
98,bottles,of,beer
90,bottles,of,beer
89,bottles,of,beer
93,bottles,of,beer
97,bottles,of,beer
96,bottles,of,beer
95,bottles,of,beer
94,bottles,of,beer
99,bottles,of,beer
92,bottles,of,beer
...
```

## License

[MIT](https://github.com/tsak/concurrent-csv-writer/blob/master/LICENSE)