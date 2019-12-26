package main

import (
	"fmt"
)

func main() {
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for index, number := range arr {
		fmt.Println(isPrime(index, number))
	}
}

func isPrime(index int, number int) (int, bool) {
	for i := 2; i < number; i++ {
		if number%i == 0 {
			return number, false
		}
	}

	return number, true
}
