package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

/**
**/

//struct to load json data into(each job)
//job might be a list of tasks
type PageVisit struct {
	ID          string `json:"id"`
	Page        string `json:"page"`
	SessionHash string `json:"sessionHash"`
}

//Jobs struct is separate the map
type Job struct {
	Date   string
	Visits []PageVisit
}

type Result struct {
	Date   string         `json:"date"`
	ByPage map[string]int `json:"byPage"`
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func allocateJobs(i chan Job, DailyStats map[string][]PageVisit) {
	for date, visits := range DailyStats {
		i <- Job{
			Date:   date,
			Visits: visits,
		}
	}
	close(i)
}

func worker(i chan Job, o chan Result, wg *sync.WaitGroup) {
	j := <-i
	vc := make(map[string]int)

	for _, v := range j.Visits {
		vc[v.Page] = vc[v.Page] + 1
	}

	o <- Result{
		Date:   j.Date,
		ByPage: vc,
	}
	wg.Done()
}

func createWorkerPool(numOfWorkers int, i chan Job, o chan Result) {
	var wg sync.WaitGroup
	for k := 0; k < numOfWorkers; k++ {
		wg.Add(1)
		go worker(i, o, &wg)
	}
	wg.Wait()
	close(o)
}

func fileResult(o chan Result) {
	var fr []Result
	for r := range o {
		fr = append(fr, r)
	}

	j, err := json.Marshal(fr)
	handleErr(err)

	err = ioutil.WriteFile("ex004/result.json", j, 0644)
	handleErr(err)
}

func main() {
	//load json into memory
	j, err := ioutil.ReadFile("ex004/data.json")
	handleErr(err)

	//load json into structs
	DailyStats := make(map[string][]PageVisit)
	err = json.Unmarshal(j, &DailyStats)
	handleErr(err)

	numberOfJobs := len(DailyStats)
	numberOfWorkers := len(DailyStats)

	//channels of communication
	jobsChan := make(chan Job, numberOfJobs)
	resultChan := make(chan Result, numberOfJobs)

	allocateJobs(jobsChan, DailyStats)
	createWorkerPool(numberOfWorkers, jobsChan, resultChan)
	fileResult(resultChan)

}
