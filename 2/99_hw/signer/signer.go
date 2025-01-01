package main

import (
	"fmt";
	"sort";
	"strings";
	"sync"
)


func ExecutePipeline(jobsArray ... job){
	var wg sync.WaitGroup
	in := make(chan interface{})
	for _ , oneJob := range jobsArray {
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


func SingleHash(in, out chan interface{}){
	for str := range in {
		crc32Chan := make(chan string)
        go func() {
            crc32Chan <- DataSignerCrc32(fmt.Sprint(str))
        }()
		
		fmt.Printf("%v SingleHash data %v\n", str, str)
		
		resultSH := DataSignerMd5(fmt.Sprint(str))
		fmt.Printf("%v SingleHash md5(data) %v\n", str, resultSH)
		
		resultSH = DataSignerCrc32(fmt.Sprint(resultSH))
		fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", str, resultSH)
		
		final := resultSH
		
        crc32 := <-crc32Chan
		fmt.Printf("%v SingleHash crc32(data) %v\n", str, crc32)
		final = crc32 + "~" + final

		fmt.Printf("%v SingleHash result %v\n", str, final)
		
		out <- final
	}
}


func MultiHash(in, out chan interface{}){
	var resultArray [6]string
	for str := range in {
		for i := 0; i< 6; i++{
			go func() {
				resultSH := DataSignerCrc32(fmt.Sprint(i) + fmt.Sprint(str))
				fmt.Printf("%v MultiHash crc32(th+step1) %v %v\n", fmt.Sprint(str), i, resultSH)
			}()
		}
		fmt.Println(resultArray)
		out <- str
	}
}


func CombineResults(in, out chan interface{}){
	var results []string

    for result := range in {
        results = append(results, fmt.Sprintf("%v", result))
    }

    sort.Strings(results)
    combined := strings.Join(results, "_")

    out <- combined
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
			fmt.Println("Вот и все" + data)
		}),
	}
	ExecutePipeline(hashSignJobs...)
}