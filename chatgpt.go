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
		out := make(chan interface{}, 100) // Буферизированный канал для увеличения производительности
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
	var mu sync.Mutex
	for str := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()

			crc32Chan := make(chan string)
			md5Chan := make(chan string)

			// Параллельно вычисляем Crc32 и Md5
			go func() {
				crc32Chan <- DataSignerCrc32(fmt.Sprint(data))
			}()
			go func() {
				mu.Lock()
				md5Chan <- DataSignerMd5(fmt.Sprint(data))
				mu.Unlock()
			}()

			md5Hash := <-md5Chan
			crc32Hash := <-crc32Chan

			crc32md5Hash := DataSignerCrc32(md5Hash)
			result := crc32Hash + "~" + crc32md5Hash
			fmt.Printf("%v SingleHash result %v\n", data, result)

			out <- result
		}(str)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for str := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()

			var resultArray [6]string
			var secondWg sync.WaitGroup

			for i := 0; i < 6; i++ {
				secondWg.Add(1)
				go func(i int) {
					defer secondWg.Done()
					resultArray[i] = DataSignerCrc32(fmt.Sprintf("%d%v", i, data))
				}(i)
			}
			secondWg.Wait()

			result := strings.Join(resultArray[:], "")
			out <- result
		}(str)
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

func main(){
	inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			fmt.Println(data)
		}),
	}
	ExecutePipeline(hashSignJobs...)
}
