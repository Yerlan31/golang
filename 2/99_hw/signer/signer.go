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
	// сюда писать код
}