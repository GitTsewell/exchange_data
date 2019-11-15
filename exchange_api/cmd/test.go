package main

import (
	"exchange_api/db"
	"fmt"
)

func main() {
	redis := db.InitRedis()
	fmt.Println(redis)
}
