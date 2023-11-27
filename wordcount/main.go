package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
)

var mu sync.Mutex

func countWords(filename string) map[string]int {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	counts := make(map[string]int)

	for scanner.Scan() {
		t := scanner.Text()
		go func(t string) {
			mu.Lock()
			if _, exist := counts[t]; !exist {
				counts[t] = 1
			} else {
				counts[t] += 1
			}
			mu.Unlock()
		}(t)
	}

	return counts
}

func main() {
	filename := os.Args

	counts := countWords(filename[1])
	keys := make([]string, 0, len(counts))

	for key := range counts {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})

	fmt.Println("Word Count:")
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, counts[k])
	}
}
