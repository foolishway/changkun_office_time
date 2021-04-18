package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

const STATUS_PATH string = "https://office.changkun.de/status"

func main() {
	var officeTime int
	var preStatus string
	var startDate int = time.Now().Day()
	var preDate int
	var officePreTime time.Time

	ticker := time.Tick(10 * time.Second)

	for range ticker {
		curDate := time.Now().Day()
		// start by tomorrow
		if startDate == curDate {
			continue
		}

		if preDate < curDate {
			f, err := os.OpenFile("./office_time", os.O_CREATE|os.O_RDWR, 0755)
			defer f.Close()
			if err != nil {
				panic(err)
			}
			msg := fmt.Sprintf("%d: %d", preDate, officeTime)
			fmt.Fprintf(f, msg)

			officeTime = 0
			preStatus = ""
		}

		var wg sync.WaitGroup
		var b []byte
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := http.Get(STATUS_PATH)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			b, err = ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
		}()
		wg.Wait()

		if string(b) == "yes!" {
			if preStatus == "yes!" {
				// calculate office time
				officeTime += int(time.Since(officePreTime))
			}
			officePreTime = time.Now()
		}

		// set office status
		preStatus = string(b)
	}
}
