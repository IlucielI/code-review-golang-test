package main

import (
	"fmt"
	"sync"
)

// Vulnerable: wg.Add(1) inside goroutine causes race condition with wg.Wait()
func fetchAllDataRace(urls []string) {
	var wg sync.WaitGroup
	for _, u := range urls {
		go func(url string) {
			wg.Add(1) // BUG: race condition! wg.Wait() in caller might return before this goroutine starts
			defer wg.Done()
			fmt.Println("fetching:", url)
		}(u)
	}
	wg.Wait()
}
