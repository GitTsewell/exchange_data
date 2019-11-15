package tool

import (
	"encoding/json"
	"exchange_api/db"
	"exchange_api/model"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type BinanceWs struct {
	*Wsbuilder
	sync.Once
	wsConn *WsConn
	redis *redis.Client

	depthData *model.BinanceDepth

	depthCallback func([]byte)

}

func NewBinanceWs() *BinanceWs {
	return &BinanceWs{}
}

func (baws *BinanceWs) BinanceSetCallback(f func([]byte)) {
	baws.depthCallback = f
}

func (baws *BinanceWs) BinanceConnect() {
	baws.Once.Do(func() {
		baws.wsConn = baws.Wsbuilder.Build()
		baws.wsConn.ReceiveMessage()
	})
}

func (baws *BinanceWs) BinanceSubscribeDepth(path string)  {
	baws.Wsbuilder = NewBuilder()
	baws.Wsbuilder.SetUrl("wss://stream.binance.com:9443"+path) // url配置
	baws.SetCheckStatusTime(time.Second * 30)
	baws.WsConfig.Handle = baws.handle
	baws.redis = db.InitRedis()
	baws.BinanceConnect()
}

func (baws *BinanceWs) handle(msg []byte)  {
	baws.wsConn.UpdateActiveTime()
	baws.depthCallback(msg)
}

func (baws *BinanceWs) BinanceDepth(msg []byte)  {
	err := json.Unmarshal(msg,&baws.depthData)
	if err == nil && len(baws.depthData.Data.Bids) > 0 {
		baws.depthToDb()
	}
}

func (baws *BinanceWs) depthToDb ()  {
	asks := ArrInterfaceToFloat64(baws.depthData.Data.Asks)
	bids := ArrInterfaceToFloat64(baws.depthData.Data.Bids)
	origin := map[string]interface{}{
		"sell":asks,
		"buy":bids,
	}
	st ,_ := json.Marshal(origin)

	rst := map[string]interface{}{
		"average_buy" : asks[0][0],
		"average_sell" : bids[0][0],
		"average_price" : asks[0][0],
		"microtime" : time.Now().UnixNano() / 1e6,
		"origin" : st,
	}

	symbol := string(baws.depthData.Stream[:len(baws.depthData.Stream)-8])
	key := fmt.Sprintf("binance:depth:0:%s",symbol)

	baws.redis.HMSet(key,rst)
}

func (baws *BinanceWs) BinanceDepthTmp(msg []byte)  {
	err := json.Unmarshal(msg,&baws.depthData)
	if err == nil && len(baws.depthData.Data.Bids) > 0 {
		symbol := string(baws.depthData.Stream[:len(baws.depthData.Stream)-8])
		if res ,_ := baws.redis.SIsMember("tmp:depth:binance",symbol).Result(); !res {
			baws.redis.SAdd("tmp:depth:binance",symbol)
		}
	}
}
