package main

import (
	"fmt";
	"sort";
	"strings"
)


func ExecutePipeline(jobsArray ... job){
	fmt.Println("ExecutePipeline!")
	in := make(chan interface{}, 1)
	for _ , oneJob := range jobsArray {
		// go oneJob(in,out)
		out := make(chan interface{}, 1)
		go oneJob(in,out)
		in = out
	}
}


func SingleHash(in, out chan interface{}){
	for str := range in {
		fmt.Printf("%v SingleHash data %v\n", str, str)

		resultSH := DataSignerMd5(fmt.Sprint(str))
		fmt.Printf("%v SingleHash md5(data) %v\n", str, resultSH)

		resultSH = DataSignerCrc32(fmt.Sprint(resultSH))
		fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", str, resultSH)
		
		final := resultSH

		resultSH = DataSignerCrc32(fmt.Sprint(str))
		fmt.Printf("%v SingleHash crc32(data) %v\n", str, resultSH)
		final = resultSH + "~" + final

		fmt.Printf("%v SingleHash result %v\n", str, final)
		
		out <- final
	}
	close(out)
}


func MultiHash(in, out chan interface{}){
	for str := range in {
		for i := 0; i< 6; i++{
			resultSH := DataSignerCrc32(fmt.Sprint(i) + fmt.Sprint(str))
			fmt.Printf("%v MultiHash crc32(th+step1) %v %v\n", fmt.Sprint(str), i, resultSH)
		}
		fmt.Printf("Закончили \n")
		out <- str
	}
	close(out)
}


func CombineResults(in, out chan interface{}){
	var results []string

    // Считываем все данные из входного канала
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

	fmt.Scanln()
}