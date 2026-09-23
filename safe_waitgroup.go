package main

import (
	"fmt"
	"sync"
)

// Safe: wg.Add(1) called synchronously before spawning goroutine
func fetchAllDataSafe(urls []string) {
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1) // Safe
		go func(url string) {
			defer wg.Done()
			fmt.Println("fetching:", url)
		}(u)
	}
	wg.Wait()
}
