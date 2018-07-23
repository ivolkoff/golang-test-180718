package main

import (
	"fmt"
	"os"
	"bufio"
	"sync"
	"net/http"
	"io/ioutil"
	"log"
	"io"
	"bytes"
)

func main() {
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	jobs := make(chan string, 5)
	var sum uint64 = 0

	//Запускаем 5 горутин-воркеров
	for w := 0; w < 5; w++ {
		go worker(jobs, wg, mutex, &sum)
	}

	reader := bufio.NewReader(os.Stdin)
	for true {
		//Читаем построчно
		text, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Panicf("Error ReadString: %s", err)
		}
		text = text[0 : len(text)-1]
		wg.Add(1)
		jobs <- text
	}
	close(jobs)

	wg.Wait()

	fmt.Printf("Total: %d", sum)
}

func worker(jobs <-chan string, wg *sync.WaitGroup, mutex *sync.Mutex, sum *uint64) {
	for url := range jobs {
		//Находим вхождение строки
		s, err := getCount(url)
		if err != nil {
			log.Panicf("Error url getCount: %s", err)
		}
		fmt.Printf("Count for %s: %d\n", url, s)
		mutex.Lock()
		*sum += uint64(s)
		mutex.Unlock()
		wg.Done()
	}
}

func getCount(url string) (i int, err error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	i = bytes.Count(buf, []byte("Go"))

	return i, nil
}
