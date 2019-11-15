package tool

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type WsConfig struct {
	url string
	HeartBeatTime time.Duration				// 心跳检测间隔时间
	heartBeatData []byte					// 心跳数据
	CheckStatusTime time.Duration			// 检测连接状态间隔时间
	ReconnectIntervalTime time.Duration   	// 定时重连时间间隔

	Handle func([]byte)
}

type WsConn struct {
	*websocket.Conn
	sync.Mutex
	WsConfig

	ActiveTime time.Time
	ActiveTimeL sync.Mutex

	subs []string
}

type Wsbuilder struct {
	WsConfig *WsConfig
}

func NewBuilder() *Wsbuilder {
	return &Wsbuilder{&WsConfig{}}
}

func (b *Wsbuilder ) SetUrl(url string) {
	b.WsConfig.url = url
}

func (b *Wsbuilder) SetHeartBeatData(data []byte) {
	b.WsConfig.heartBeatData = data
}

func (b *Wsbuilder) SetHeartBeatTime(time time.Duration) {
	b.WsConfig.HeartBeatTime = time
}

func (b *Wsbuilder) SetCheckStatusTime(time time.Duration) {
	b.WsConfig.CheckStatusTime = time
}
func (b *Wsbuilder) Build() *WsConn {
	WsConn := &WsConn{WsConfig:*b.WsConfig}
	return WsConn.NewWs()
}

func (ws *WsConn) NewWs() *WsConn {
	ws.Lock()
	defer ws.Unlock()

	if err := ws.connect();err != nil {
		log.Println(ws.url, "ws connect error ", err)
	}

	ws.ActiveTime = time.Now()
	ws.heartBeat()
	ws.checkStatusTimer()

	return ws
}

func (ws *WsConn) checkStatusTimer()  {
	checkStatusTimer := time.NewTicker(ws.CheckStatusTime)

	go func() {
		for {
			select {
			case <-checkStatusTimer.C:
				if time.Now().Sub(ws.ActiveTime) > ws.CheckStatusTime *2 {
					ws.ReConnect()
				}
			}
		}
	}()
}

func (ws *WsConn) heartBeat()  {
	if ws.HeartBeatTime == 0 {
		return
	}

	wsHeart := time.NewTicker(ws.HeartBeatTime)

	go func() {
		for  {
			select {
			case <-wsHeart.C:
				_ = ws.SendMessage(ws.heartBeatData)
				ws.UpdateActiveTime()
			}
		}
	}()
}

func (ws *WsConn) UpdateActiveTime()  {
	ws.ActiveTimeL.Lock()
	defer ws.ActiveTimeL.Unlock()

	ws.ActiveTime = time.Now()
}

func (ws *WsConn) connect() error {
	log.Printf("connecting to %s", ws.url)
	conn ,_,err := websocket.DefaultDialer.Dial(ws.url,nil)
	if err != nil {
		return err
	}
	ws.Conn = conn

	return nil
}

func (ws *WsConn) ReConnect() {
	ws.Lock()
	defer ws.Unlock()

	log.Println("close ws  error :", ws.Close())
	time.Sleep(time.Second)

	if err := ws.connect(); err != nil {
		log.Println(ws.url, "ws connect error ", err)
		return
	}

	for _,sub := range ws.subs {
		log.Println("subscribe:", sub)
		_ = ws.SendMessage([]byte(sub))
	}

}

func (ws *WsConn) Subscribe(subEvent string) error {
	log.Println("Subscribe:", subEvent)

	err := ws.SendMessage([]byte(subEvent))
	if err != nil {
		return err
	}
	ws.subs = append(ws.subs, subEvent)
	return nil
}

func (ws *WsConn) SendMessage(msg []byte) error {
	ws.Lock()
	defer ws.Unlock()

	err := ws.WriteMessage(websocket.TextMessage,msg)
	if err != nil {
		return err
	}

	return nil
}

func (ws *WsConn) ReceiveMessage() {
	go func() {
		for {
			_,message ,err := ws.ReadMessage()

			if err != nil {
				log.Println("websocket消息读取失败 : ",err)
			}
			ws.Handle([]byte(message))
		}
	}()
}