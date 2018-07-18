package main

import (
	"fmt"
	"os"
	"bufio"
	"sync"
	"strings"
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
		text, _ := reader.ReadString('\n')
		if len(text) == 0 {
			break
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
	for str := range jobs {
		//Находим вхождение строки
		i := strings.Index(str, "go")
		fmt.Printf("Count for %s: %d\n", str, i)
		mutex.Lock()
		*sum += uint64(i)
		mutex.Unlock()
		wg.Done()
	}
}
