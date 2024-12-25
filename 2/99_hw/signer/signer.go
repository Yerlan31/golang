package main

import (
	"fmt"
)

func ExecutePipeline(jobsArray ... job){
	fmt.Println("ExecutePipeline!")
	in := make(chan interface{})
    out := make(chan interface{})
	for _ , oneJob := range jobsArray {
		go oneJob(in,out)
	}
}

func SingleHash(in, out chan interface{}){
	fmt.Println("SingleHash!")
}

func MultiHash(in, out chan interface{}){
	fmt.Println("MultiHash!")
}

func CombineResults(in, out chan interface{}){
	fmt.Println("CombineResults!")
}

func main(){
	fmt.Println("Hello!")
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

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