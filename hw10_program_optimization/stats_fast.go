package hw10programoptimization

import (
	"bufio"
	"io"
	"runtime"
	"strings"
	"sync"

	"github.com/valyala/fastjson"
)

func GetDomainStatFastButGreedy(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	cpuCount := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuCount)
	stat := make(DomainStat)

	wg := new(sync.WaitGroup)
	wg1 := new(sync.WaitGroup)

	lock := sync.Mutex{}

	jobs := make(chan string, cpuCount)
	results := make(chan string, cpuCount)

	for i := 0; i < cpuCount; i++ {
		wg.Add(1)
		go worker(wg, wg1, jobs, results, domain)
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

func worker(wg, wg1 *sync.WaitGroup, jobs <-chan string, results chan<- string, domain string) {
	for line := range jobs {
		s, err := job(line, domain)
		if err != nil {
			continue
		}
		wg1.Add(1)
		results <- s
	}
	wg.Done()
}

func job(line, domain string) (string, error) {
	email := fastjson.GetString([]byte(line), "Email")

	if !strings.Contains(email, "@") {
		return "", ErrInvalidEmail
	}

	if strings.HasSuffix(email, "."+domain) {
		return strings.ToLower(strings.SplitN(email, "@", 2)[1]), nil
	}

	return "", ErrDomainNotFound
}
