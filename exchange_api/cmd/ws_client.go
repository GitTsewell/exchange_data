package main

import (
	"exchange_api/db"
	"exchange_api/tool"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

func main()  {
	redis := db.InitRedis()
	defer redis.Close()

	okexRun ,_ := redis.Get("config:exchange:okex").Result()
	bitmexRun ,_ := redis.Get("config:exchange:bitmex").Result()
	huobiRun ,_ := redis.Get("config:exchange:huobi").Result()
	binanceRun ,_ := redis.Get("config:exchange:binance").Result()

	// okex
	if okexRun != "0" {
		go okexWs(redis)
	}

	// bitmex
	if bitmexRun != "0" {
		go bitmexWs(redis)
	}

	//火币现货
	if huobiRun != "0" {
		go huobiSpotWs(redis)
	}

	// 火币期货
	if huobiRun != "0" {
		go huobiFutureWs(redis)
	}

	// binance
	if binanceRun != "0" {
		go binanceWs(redis)
	}

	for {
		select {}
	}

}

func okexWs(redis *redis.Client)  {
	r1,_ := redis.SMembers("config:depth:okex:spot").Result()
	r2,_ := redis.SMembers("config:depth:okex:future").Result()

	if len(r1) > 0 || len(r2) > 0 {
		ws := tool.NewOkexWs()
		ws.OkexSetCallback(func(msg []byte) {
			ws.OkexDepth(msg)
		})

		for _,v1 := range r1{
			sub := fmt.Sprintf(`{"op":"subscribe","args":["spot/depth5:%s"]}`,v1)
			ws.OkexSubscribeDepth(sub)
		}

		for _,v2 := range r2 {
			sub := fmt.Sprintf(`{"op":"subscribe","args":["futures/depth5:%s"]}`,v2)
			ws.OkexSubscribeDepth(sub)
		}
	}
}

func bitmexWs(redis *redis.Client)  {
	r,_ := redis.SMembers("config:depth:bitmex:future").Result()

	if len(r) > 0 {
		ws := tool.NewBitmexWs()
		ws.BitmexSetCallback(func(msg []byte) {
			ws.BitmexDepth(msg)
		})

		for _,v := range r{
			sub := fmt.Sprintf(`{"op": "subscribe", "args":["orderBook10:%s"]}`,v)
			ws.BitmexSubscribeDepth(sub)
		}
	}
}

func huobiSpotWs(redis *redis.Client)  {
	r,_ := redis.SMembers("config:depth:huobi:spot").Result()

	if len(r) > 0 {
		ws := tool.NewHuobiWs("wss://api.huobi.pro/ws")
		ws.HuobiSetCallback(func(msg []byte) {
			ws.HuobiDepth(msg)
		})

		for i,v := range r{
			sub := fmt.Sprintf(`{"id":"id_%s","sub":"market.%s.depth.step0"}`,i,v)
			ws.HuobiSubscribeDepth(sub)
		}
	}
}

func huobiFutureWs(redis *redis.Client)  {
	r,_ := redis.SMembers("config:depth:huobi:future").Result()

	if len(r) > 0 {
		ws := tool.NewHuobiWs("wss://www.hbdm.com/ws")
		ws.HuobiSetCallback(func(msg []byte) {
			ws.HuobiDepth(msg)
		})

		for i,v := range r{
			sub := fmt.Sprintf(`{"id":"id_%s","sub":"market.%s.depth.step0"}`,i,v)
			ws.HuobiSubscribeDepth(sub)
		}
	}
}

func binanceWs(redis *redis.Client)  {
	r,_ := redis.SMembers("config:depth:binance:spot").Result()

	if len(r) > 0 {
		ws := tool.NewBinanceWs()
		ws.BinanceSetCallback(func(msg []byte) {
			ws.BinanceDepth(msg)
		})

		split := []string{}
		for _,i := range r{
			split = append(split,i+"@depth20")
		}
		stream := strings.Join(split,"/")
		keys := fmt.Sprintf("/stream?streams=%s",stream)

		ws.BinanceSubscribeDepth(keys)
	}
}
