package main

import (
	"runtime"
	"sync"
)

// DispatchConcurrencyWorkers launches cpu number of concurrency workers to execute proc function with targets
func DispatchConcurrencyWorkers(targets []string, proc func(string) XlsxMap) XlsxMap {
	size := len(targets)
	targetsChan := make(chan string, size)
	for _, target := range targets {
		targetsChan <- target
	}

	var wg sync.WaitGroup
	wg.Add(size)

	out := make(chan XlsxMap, size)
	for i := 0; i < runtime.NumCPU(); i++ {
		go LaunchWorker(targetsChan, out, &wg, proc)
	}

	close(targetsChan)
	wg.Wait()

	close(out)

	return MergeWorkerResults(out)
}

// LaunchWorker executes proc function until targets channel is closed
func LaunchWorker(targets chan string, out chan XlsxMap, wg *sync.WaitGroup, proc func(string) XlsxMap) {
	for target := range targets {
		out <- proc(target)
		wg.Done()
	}
}

// MergeWorkerResults merges some XlsxMap from out channel into one XlsxMap
func MergeWorkerResults(out chan XlsxMap) XlsxMap {
	ret := XlsxMap{}
	for parsed := range out {
		for k, v := range parsed {
			ret[k] = v
		}
	}
	return ret
}
