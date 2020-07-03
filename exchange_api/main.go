package main

import (
	"exchange_api/cfg"
	"exchange_api/route"
	"fmt"
	"os"
)

func main() {
	// init cfg
	err := cfg.Initialize("./application")
	if err != nil {
		fmt.Printf("cfg err :%s", err.Error())
		os.Exit(1)
	}

	r := route.InitRoute()
	_ = r.Run()
}
