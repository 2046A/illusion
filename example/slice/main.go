package main

import "fmt"

func main() {
	base := [5]int{1, 2, 3, 4, 5}
	slice := make([]int, 0, 10)
	slice = append(slice, base[0], base[1])
	for _, i := range slice {
		fmt.Println(i)
	}
	fmt.Println("******************")
	slice = slice[0:0]
	slice = append(slice, 123)
	for _, i := range slice {
		fmt.Println(i)
	}
}
