package main

import "fmt"

func main() {
	maxSize := 50
	maxArray := make([]int, maxSize)
	nVars := 5
	var s []int = maxArray[1:nVars]
	fmt.Println(s)
}
