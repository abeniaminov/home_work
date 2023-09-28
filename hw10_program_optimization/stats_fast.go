package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync"
)



func GetDomainStatFast(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	cpuCount := runtime.NumCPU()
	stat := make(DomainStat)
		
	reg, err := regexp.Compile("\\."+domain)
	if err != nil {
		return stat, err
	}

    var wg = new(sync.WaitGroup)
	var wg1 = new(sync.WaitGroup)

	var lock = sync.RWMutex{}

	jobs := make(chan string, cpuCount)
	results := make(chan string, cpuCount)

	for i := 0; i < cpuCount; i++ {
		wg.Add(1)
		go worker(wg, wg1, jobs, results, reg)
	}

	go func() {
		for scanner.Scan() {
			jobs <- scanner.Text()
		}
		close(jobs)
	}()

    go func() {
		for result := range results {
			lock.Lock()
			stat[result]++
			lock.Unlock()
			wg1.Done() 
		}
	}()

	wg.Wait()
	close(results)
	wg1.Wait()

	return stat, nil	
}

func worker(wg, wg1 *sync.WaitGroup, jobs <-chan string, results chan<- string, reg *regexp.Regexp) {
    for line := range jobs {
		s, err := job(line, reg)
		if err != nil {
			continue
		}
		wg1.Add(1)
        results <- s

    }
   wg.Done()
}

func job(line string, reg *regexp.Regexp) (string, error) {
	var user User
	if err := json.Unmarshal([]byte(line), &user); err != nil {
		return "", err
	}
	
    if reg.Match([]byte(user.Email)) {
		return strings.ToLower(strings.SplitN(user.Email, "@", 2)[1]), nil
	}
	
	return "", ErrDomainNotFound 	
}

	

