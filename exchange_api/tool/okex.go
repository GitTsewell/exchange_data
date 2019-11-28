package tool

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"exchange_api/config"
	"exchange_api/db"
	"exchange_api/model"
	"fmt"
	"github.com/go-redis/redis"
	"io/ioutil"
	"sync"
	"time"
)

type OkexWs struct {
	*Wsbuilder
	sync.Once
	wsConn *WsConn
	redis *redis.Client

	depthData *model.OkexDepth

	depthCallback func([]byte)

}

func NewOkexWs() *OkexWs {
	okexWs := &OkexWs{Wsbuilder:NewBuilder()}
	okexWs.redis = db.InitRedis()
	okexWs.Wsbuilder.SetUrl("wss://real.okex.com:8443/ws/v3") // url配置
	okexWs.Wsbuilder.SetHeartBeatData([]byte("ping")) // ping
	okexWs.Wsbuilder.SetHeartBeatTime(time.Second * 30) // 检测时间
	okexWs.Wsbuilder.SetCheckStatusTime(time.Second * 30) // 检测时间
	okexWs.WsConfig.Handle = okexWs.handle
	return okexWs
}

func (okws *OkexWs) OkexSetCallback(f func([]byte)) {
	okws.depthCallback = f
}

func (okws *OkexWs) OkexConnect() {
	okws.Once.Do(func() {
		okws.wsConn = okws.Wsbuilder.Build()
		okws.wsConn.ReceiveMessage()
	})
}

func (okws *OkexWs) OkexSubscribeDepth(msg string)  {
	okws.OkexConnect()
	_ = okws.wsConn.Subscribe(msg)
}

func (okws *OkexWs) handle(msg []byte)  {
	reader :=flate.NewReader(bytes.NewBuffer(msg))
	defer reader.Close()

	text,_ := ioutil.ReadAll(reader)

	if string(text) == "pong" {
		okws.wsConn.UpdateActiveTime()
	}

	okws.depthCallback(text)
}

func (okws *OkexWs) OkexDepth(msg []byte)  {
	err := json.Unmarshal(msg,&okws.depthData)
	if err == nil && len(okws.depthData.Data) > 0 {
		okws.depthToDb()
	}
}

func (okws *OkexWs) depthToDb ()  {
	asks := ArrInterfaceToFloat64(okws.depthData.Data[0].Asks)
	bids := ArrInterfaceToFloat64(okws.depthData.Data[0].Bids)
	origin := map[string]interface{}{
		"sell":asks,
		"buy":bids,
	}
	st ,_ := json.Marshal(origin)

	rst := map[string]interface{}{
		"average_buy" : okws.depthData.Data[0].Asks[0][0],
		"average_sell" : okws.depthData.Data[0].Bids[0][0],
		"average_price" : okws.depthData.Data[0].Asks[0][0],
		"microtime" : time.Now().UnixNano() / 1e6,
		"origin" : st,
	}

	var key string

	if okws.depthData.Table == config.OKEX_DEPTH_TABLE_SPOT {
		key = fmt.Sprintf("okex:depth:%s:%s",config.SPOT,okws.depthData.Data[0].InstrumentID)
	}else {
		key = fmt.Sprintf("okex:depth:%s:%s",config.FUTURE,okws.depthData.Data[0].InstrumentID)
	}

	okws.redis.HMSet(key,rst)
	okws.redis.Expire(key,time.Minute * 5)
}

func (okws *OkexWs) OkexDepthTmp(msg []byte)  {
	err := json.Unmarshal(msg,&okws.depthData)
	if err == nil && len(okws.depthData.Data) > 0 {
		if res ,_ := okws.redis.SIsMember("tmp:depth:okex",okws.depthData.Data[0].InstrumentID).Result(); !res {
			okws.redis.SAdd("tmp:depth:okex",okws.depthData.Data[0].InstrumentID)
		}

	}
}