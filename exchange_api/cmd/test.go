package main

import "fmt"

func main() {
	arr := [8]int{0,1,2,3,4,5,6,7}

	spl := arr[2:5]

	fmt.Println(spl)
	fmt.Println(len(spl))
	fmt.Println(cap(spl))

}
