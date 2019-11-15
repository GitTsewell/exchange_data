package main

import "exchange_api/route"

func main() {
	r := route.InitRoute()
	_ = r.Run()
}
