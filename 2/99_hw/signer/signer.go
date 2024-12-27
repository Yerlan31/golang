package main

import (
	"fmt"
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
	str := <-in
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


func MultiHash(in, out chan interface{}){
	str := <-in
	for i := 0; i< 6; i++{
		resultSH := DataSignerCrc32(fmt.Sprint(i) + fmt.Sprint(str))
		fmt.Printf("%v MultiHash crc32(th+step1) %v\n", str, i, resultSH)
	}
	out <- str
}


func CombineResults(in, out chan interface{}){
	<-in
	fmt.Println("CombineResults!")
	out <- 3
}


func main(){
	inputData := []int{0, 1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				fmt.Println("Первая функция")
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

	fmt.Scanln()
}