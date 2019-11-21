package main

import (
	"exchange_api/db"
	"fmt"
)

func main() {
	redis := db.InitRedis()
	defer redis.Close()

	if a,_ :=redis.Set("asdasd","asdasd",0).Result();a == "OK" {
		fmt.Println(111)
	}else {
		fmt.Println(222)
	}

	a := []map[string]interface{}{}
	


}
