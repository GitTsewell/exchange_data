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

type BitmexWs struct {
	*Wsbuilder
	sync.Once
	wsConn *WsConn
	redis *redis.Client

	depthData *model.BitmexDepth

	depthCallback func([]byte)

}

func NewBitmexWs() *BitmexWs {
	bitmexWs := &BitmexWs{Wsbuilder:NewBuilder()}
	bitmexWs.redis = db.InitRedis()
	bitmexWs.Wsbuilder.SetUrl("wss://www.bitmex.com/realtime") // url配置
	bitmexWs.Wsbuilder.SetHeartBeatData([]byte("ping")) // ping
	bitmexWs.Wsbuilder.SetHeartBeatTime(time.Second * 30) // 检测时间
	bitmexWs.Wsbuilder.SetCheckStatusTime(time.Second * 30) // 状态检测时间
	bitmexWs.WsConfig.Handle = bitmexWs.handle
	return bitmexWs
}

func (btws *BitmexWs) BitmexSetCallback(f func([]byte)) {
	btws.depthCallback = f
}

func (btws *BitmexWs) BitmexConnect() {
	btws.Once.Do(func() {
		btws.wsConn = btws.Wsbuilder.Build()
		btws.wsConn.ReceiveMessage()
	})
}

func (btws *BitmexWs) BitmexSubscribeDepth(msg string)  {
	btws.BitmexConnect()
	_ = btws.wsConn.Subscribe(msg)
}

func (btws *BitmexWs) handle(msg []byte)  {

	if string(msg) == "pong" {
		btws.wsConn.UpdateActiveTime()
	}

	btws.depthCallback(msg)
}

func (btws *BitmexWs) BitmexDepth(msg []byte)  {
	err := json.Unmarshal(msg,&btws.depthData)
	if err == nil && len(btws.depthData.Data) > 0 {
		btws.depthToDb()
	}
}

func (btws *BitmexWs) depthToDb ()  {
	origin := map[string]interface{}{
		"sell":btws.depthData.Data[0].Asks,
		"buy":btws.depthData.Data[0].Bids,
	}
	st ,_ := json.Marshal(origin)

	rst := map[string]interface{}{
		"average_buy" : btws.depthData.Data[0].Asks[0][0],
		"average_sell" : btws.depthData.Data[0].Bids[0][0],
		"average_price" : btws.depthData.Data[0].Asks[0][0],
		"microtime" : time.Now().UnixNano() / 1e6,
		"origin" : st,
	}

	key := fmt.Sprintf("bitmex:depth:1:%s",btws.depthData.Data[0].Symbol)

	btws.redis.HMSet(key,rst)
}

func (btws *BitmexWs) BitmexDepthTmp(msg []byte)  {
	err := json.Unmarshal(msg,&btws.depthData)
	if err == nil && len(btws.depthData.Data) > 0 {
		if res ,_ := btws.redis.SIsMember("tmp:depth:bitmex",btws.depthData.Data[0].Symbol).Result(); !res {
			btws.redis.SAdd("tmp:depth:bitmex",btws.depthData.Data[0].Symbol)
		}
	}
}