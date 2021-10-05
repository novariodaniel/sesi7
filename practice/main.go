package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConnString   = "username:password@tcp(127.0.0.1:3306)/kampus-merdeka?charset=utf8mb4"
	dbMaxIdleConns = 4
	dbMaxConns     = 100
	totalWorker    = 100
	csvFile        = "majestic_million.csv"
	dataHeaders    = []string{}
)

func main() {
    start := time.Now()

    db, err := openDbConnection()
    if err != nil {
        log.Fatal(err.Error())
    }

    csvReader, csvFile, err := openCsvFile()
    if err != nil {
        log.Fatal(err.Error())
    }
    defer csvFile.Close()

    jobs := make(chan []interface{}, 0)
    wg := new(sync.WaitGroup)

    go dispatchWorkers(db, jobs, wg)
    readCsvFilePerLineThenSendToWorker(csvReader, jobs, wg)

    wg.Wait()

    duration := time.Since(start)
    fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
}

// OPEN DB CONNECTION
func openDbConnection() (*sql.DB, error) {
    log.Println("=> open db connection")

    db, err := sql.Open("mysql", dbConnString)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(dbMaxConns)
    db.SetMaxIdleConns(dbMaxIdleConns)

    return db, nil
}

// OPEN FILE
func openCsvFile() (*csv.Reader, *os.File, error) {
    log.Println("=> open csv file")

    f, err := os.Open(csvFile)
    if err != nil {
        return nil, nil, err
    }

    reader := csv.NewReader(f)
    return reader, f, nil
}

// jalankan worker
func dispatchWorkers(db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
    for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
        go func(workerIndex int, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
            counter := 0

            for job := range jobs {
                doTheJob(workerIndex, counter, db, job)
                wg.Done()
                counter++
            }
        }(workerIndex, db, jobs, wg)
    }
}

// READ FILE
func readCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- []interface{}, wg *sync.WaitGroup) {
    for {
        row, err := csvReader.Read()
        if err != nil {
            if err == io.EOF {
                err = nil
            }
            break
        }

        if len(dataHeaders) == 0 {
            dataHeaders = row
            continue
        }

        rowOrdered := make([]interface{}, 0)
        for _, each := range row {
            rowOrdered = append(rowOrdered, each)
        }

        wg.Add(1)
        jobs <- rowOrdered
    }
    close(jobs)
}

// INSERT DB
func doTheJob(workerIndex, counter int, db *sql.DB, values []interface{}) {
    for {
        var outerError error
        func(outerError *error) {
            defer func() {
                if err := recover(); err != nil {
                    *outerError = fmt.Errorf("%v", err)
                }
            }()

            conn, err := db.Conn(context.Background())
            query := fmt.Sprintf("INSERT INTO domain (%s) VALUES (%s)",
                strings.Join(dataHeaders, ","),
                strings.Join(generateQuestionsMark(len(dataHeaders)), ","),
            )

            _, err = conn.ExecContext(context.Background(), query, values...)
            if err != nil {
                log.Fatal(err.Error())
            }

            err = conn.Close()
            if err != nil {
                log.Fatal(err.Error())
            }
        }(&outerError)
        if outerError == nil {
            break
        }
    }

    if counter%100 == 0 {
        log.Println("=> worker", workerIndex, "inserted", counter, "data")
    }
}

func generateQuestionsMark(n int) []string {
    s := make([]string, 0)
    for i := 0; i < n; i++ {
        s = append(s, "?")
    }
    return s
}