package tool

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func TmpOkexDepthWs(spot,future *[]string)  {
	go func() {
		ws := NewOkexWs()
		ws.OkexSetCallback(func(msg []byte) {
			ws.OkexDepthTmp(msg)
		})

		for _,i1 := range *spot{
			key := fmt.Sprintf(`{"op":"subscribe","args":["spot/depth5:%s"]}`,i1)
			ws.OkexSubscribeDepth(key)
		}

		for _,i2 := range *future{
			key := fmt.Sprintf(`{"op":"subscribe","args":["futures/depth5:%s"]}`,i2)
			ws.OkexSubscribeDepth(key)
		}

		t := time.After(time.Second * 30)

		for  {
			select {
			case <-t:
				log.Println("okex 行情连接测试结束")
				return
			}
		}
	}()
}

func TmpHuobiDepthWs(spot,future *[]string)  {
	go func() { // 火币现货
		ws := NewHuobiWs("wss://api.huobi.pro/ws")
		ws.HuobiSetCallback(func(msg []byte) {
			ws.HuobiDepthTmp(msg)
		})

		for _,i:=range *spot{
			key := fmt.Sprintf(`{"id":"I0TwXU3P","sub":"market.%s.depth.step0"}`,i)
			ws.HuobiSubscribeDepth(key)
		}

		t := time.After(time.Second * 30)

		for  {
			select {
			case <-t:
				log.Println("Huobi 现货行情连接测试结束")
				return
			}
		}
	}()

	go func() { // 火币期货
		ws := NewHuobiWs("wss://www.hbdm.com/ws")
		ws.HuobiSetCallback(func(msg []byte) {
			ws.HuobiDepthTmp(msg)
		})

		for _,i := range *future{
			key := fmt.Sprintf(`{"id":"I0TwH86H","sub":"market.%s.depth.step0"}`,i)
			ws.HuobiSubscribeDepth(key)
		}

		t := time.After(time.Second * 30)

		for  {
			select {
			case <-t:
				log.Println("Huobi 期货行情连接测试结束")
				return
			}
		}
	}()
}

func TmpBitmexDepthWs(spot,future *[]string)  {
	go func() {
		ws := NewBitmexWs()
		ws.BitmexSetCallback(func(msg []byte) {
			ws.BitmexDepthTmp(msg)
		})

		for _,i:=range *future{
			key := fmt.Sprintf(`{"op": "subscribe", "args":["orderBook10:%s"]}`,i)
			ws.BitmexSubscribeDepth(key)
		}

		t := time.After(time.Second * 30)

		for  {
			select {
			case <-t:
				log.Println("Bitmex 行情连接测试结束")
				return
			}
		}
	}()
}

func TmpBinanceDepthWs(spot,future *[]string)  {
	go func() {
		ws := NewBinanceWs()
		ws.BinanceSetCallback(func(msg []byte) {
			ws.BinanceDepthTmp(msg)
		})
		split := []string{}
		for _,i := range *spot{
			split = append(split,i+"@depth20")
		}
		stream := strings.Join(split,"/")
		keys := fmt.Sprintf("/stream?streams=%s",stream)
		ws.BinanceSubscribeDepth(keys)

		t := time.After(time.Second * 30)

		for  {
			select {
			case <-t:
				log.Println("Binance 行情连接测试结束")
				return
			}
		}
	}()
}