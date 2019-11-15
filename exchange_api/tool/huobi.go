package tool

import (
	"encoding/json"
	"exchange_api/config"
	"exchange_api/db"
	"exchange_api/model"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HuobiWs struct {
	*Wsbuilder
	sync.Once
	wsConn *WsConn
	redis *redis.Client

	depthData *model.HuobiDepth

	depthCallback func([]byte)
}

func NewHuobiWs(url string) *HuobiWs {
	huobiWs := &HuobiWs{Wsbuilder:NewBuilder()}
	huobiWs.redis = db.InitRedis()
	huobiWs.Wsbuilder.SetUrl(url) // url配置
	huobiWs.Wsbuilder.SetCheckStatusTime(time.Second * 30) // 检测时间
	huobiWs.WsConfig.Handle = huobiWs.handle
	return huobiWs
}

func (hbws *HuobiWs) HuobiSetCallback(f func([]byte)) {
	hbws.depthCallback = f
}

func (hbws *HuobiWs) HuobiConnect() {
	hbws.Once.Do(func() {
		hbws.wsConn = hbws.Wsbuilder.Build()
		hbws.wsConn.ReceiveMessage()
	})
}

func (hbws *HuobiWs) HuobiSubscribeDepth(msg string)  {
	hbws.HuobiConnect()
	_ = hbws.wsConn.Subscribe(msg)
}

func (hbws *HuobiWs) handle(msg []byte)  {
	text,err := GzipDecodeHuobi(msg)
	if err != nil {
		fmt.Println(err)
	}

	if strings.Contains(string(text),"ping") {
		str := strconv.FormatInt(time.Now().Unix(),10)
		pong := `{"pong": ` + str + `}`
		_ = hbws.wsConn.SendMessage([]byte(pong))
		hbws.wsConn.UpdateActiveTime()
	}
	hbws.depthCallback(text)
}

func (hbws *HuobiWs) HuobiDepth (msg []byte)  {
	err := json.Unmarshal(msg,&hbws.depthData)
	if err == nil && len(hbws.depthData.Tick.Bids) > 0 {
		hbws.depthToDb()
	}
}

func (hbws *HuobiWs) depthToDb ()  {
	origin := map[string]interface{}{
		"sell":hbws.depthData.Tick.Asks,
		"buy":hbws.depthData.Tick.Bids,
	}
	st ,_ := json.Marshal(origin)

	rst := map[string]interface{}{
		"average_buy" : hbws.depthData.Tick.Asks[0][0],
		"average_sell" : hbws.depthData.Tick.Bids[0][0],
		"average_price" : hbws.depthData.Tick.Asks[0][0],
		"microtime" : time.Now().UnixNano() / 1e6,
		"origin" : st,
	}

	var key string

	chs := strings.Split(hbws.depthData.Ch,".")
	if strings.Contains(chs[1],"_") { // 期货
		key = fmt.Sprintf("huobi:depth:%s:%s",config.FUTURE,chs[1])
	}else {
		key = fmt.Sprintf("huobi:depth:%s:%s",config.SPOT,chs[1])
	}

	hbws.redis.HMSet(key,rst)
}

func (hbws *HuobiWs) HuobiDepthTmp (msg []byte)  {
	err := json.Unmarshal(msg,&hbws.depthData)
	if err == nil && len(hbws.depthData.Tick.Bids) > 0 {
		chs := strings.Split(hbws.depthData.Ch,".")
		if res,_ := hbws.redis.SIsMember("tmp:depth:huobi",chs[1]).Result(); !res {
			hbws.redis.SAdd("tmp:depth:huobi",chs[1])
		}
	}
}