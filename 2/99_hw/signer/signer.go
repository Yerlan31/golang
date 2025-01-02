package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(jobsArray ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})
	for _, oneJob := range jobsArray {
		wg.Add(1)
		out := make(chan interface{})
		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			j(in, out)
		}(oneJob, in, out)
		in = out
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	var md5Mutex sync.Mutex // Мьютекс для контроля вызовов DataSignerMd5

	for str := range in {
		wg.Add(1)

		go func(data string) {
			defer wg.Done()

			crc32Chan := make(chan string)
			crc32Md5Chan := make(chan string)

			// Вычисляем CRC32(data) параллельно
			go func() {
				crc32Chan <- DataSignerCrc32(data)
			}()

			// Вычисляем md5(data) с мьютексом, и сразу CRC32(md5(data))
			go func() {
				md5Mutex.Lock()
				md5Hash := DataSignerMd5(data)
				md5Mutex.Unlock()

				fmt.Printf("%v SingleHash md5(data) %v\n", data, md5Hash)
				crc32Md5Chan <- DataSignerCrc32(md5Hash)
			}()

			// Ждем результатов
			crc32Data := <-crc32Chan
			fmt.Printf("%v SingleHash crc32(data) %v\n", data, crc32Data)

			crc32Md5 := <-crc32Md5Chan
			fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", data, crc32Md5)

			// Формируем итоговый результат
			finalResult := crc32Data + "~" + crc32Md5
			fmt.Printf("%v SingleHash result %v\n", data, finalResult)

			out <- finalResult
		}(fmt.Sprint(str))
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for str := range in {
		wg.Add(1)

		go func(data string) {
			defer wg.Done()

			var localResults [6]string
			var innerWg sync.WaitGroup

			// Вычисляем 6 значений CRC32 параллельно
			for i := 0; i < 6; i++ {
				innerWg.Add(1)
				go func(i int) {
					defer innerWg.Done()
					th := fmt.Sprint(i)
					localResults[i] = DataSignerCrc32(th + data)
					fmt.Printf("%v MultiHash crc32(th+step1) %v %v\n", data, i, localResults[i])
				}(i)
			}

			innerWg.Wait()
			finalResult := strings.Join(localResults[:], "")
			fmt.Printf("%v MultiHash result %v\n", data, finalResult)
			out <- finalResult
		}(fmt.Sprint(str))
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []string

	for result := range in {
		results = append(results, fmt.Sprintf("%v", result))
	}

	sort.Strings(results)
	out <- strings.Join(results, "_")
}
