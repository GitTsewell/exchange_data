package app

import (
	"exchange_api/db"
	"exchange_api/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type data struct {
	Symbols     []string `form:"symbols" json:"symbols" binding:"required"`
}

func DepthIndex(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	platform := c.DefaultQuery("platform","okex")

	key := fmt.Sprintf("%s:depth:*",platform)
	keys,_ := redis.Keys(key).Result()

	now := time.Now().UnixNano() / 1e6

	list := []map[string]interface{}{}
	for _,v := range keys{
		rst ,_ := redis.HMGet(v,"average_price","microtime").Result()
		// status
		status := 1
		t ,_ := strconv.ParseInt(rst[1].(string), 10, 64)
		if now - t >= 30000 {
			status = 0
		}
		// symbol
		spli := strings.Split(v,":")

		// time
		rst[1],_ = tool.MsToTime(rst[1].(string))

		raw := map[string]interface{}{
			"symbol":spli[3],
			"price":rst[0],
			"tag":spli[2],
			"time": rst[1],
			"status":status,
		}

		list = append(list,raw)
	}

	c.JSON(200,gin.H{
		"data":list,
		"status":1,
	})

}

func DepthEdit(c *gin.Context)  {
	platform := c.DefaultQuery("platform","okex")

	redis := db.InitRedis()
	defer redis.Close()

	key1 := fmt.Sprintf("config:depth:%s:spot",platform)
	key2 := fmt.Sprintf("config:depth:%s:future",platform)
	a1,_ := redis.SMembers(key1).Result()
	a2,_ := redis.SMembers(key2).Result()

	a1 = append(a1,a2...)

	c.JSON(200,gin.H{
		"symbols":a1,
	})
}

func DepthUpdate(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	// platform
	platform := c.Param("platform")

	var data data
	_ = c.BindJSON(&data)

	key1 := fmt.Sprintf("tmp:config:depth:%s:spot",platform)
	key2 := fmt.Sprintf("tmp:config:depth:%s:future",platform)

	redis.Del(key1)
	redis.Del(key2)

	var spot []string
	var future []string

	var isSpot func(string) bool
	var tmpDepth func(*[]string,*[]string)
	switch platform {
	case "okex":
		isSpot = tool.IsSpotOfOkex
		tmpDepth = tool.TmpOkexDepthWs
		break
	case "huobi":
		isSpot = tool.IsSpotOfHuobi
		tmpDepth = tool.TmpHuobiDepthWs
		break
	case "bitmex":
		isSpot = tool.IsSpotOfBitmex
		tmpDepth = tool.TmpBitmexDepthWs
		break
	case "binance":
		isSpot = tool.IsSpotOfBinance
		tmpDepth = tool.TmpBinanceDepthWs
		break
	default:
		break
	}

	for _,i:= range data.Symbols{
		if isSpot(i) {
			spot = append(spot,i)
		}else {
			future = append(future,i)
		}
	}

	redis.SAdd(key1,spot)
	redis.SAdd(key2,future)

	tmpKey := fmt.Sprintf("tmp:depth:%s",platform)
	redis.Del(tmpKey)

	tmpDepth(&spot,&future)

	c.JSON(200,gin.H{"status":true})
}

func DepthCheck(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	platform := c.Param("platform")
	k1 := fmt.Sprintf("tmp:config:depth:%s:spot",platform)
	k2 := fmt.Sprintf("tmp:config:depth:%s:future",platform)
	k3 := fmt.Sprintf("tmp:depth:%s",platform)

	a1,_ := redis.SMembers(k1).Result()
	a2,_ := redis.SMembers(k2).Result()
	a3,_ := redis.SMembers(k3).Result()

	a1 = append(a1,a2...)

	list := []map[string]interface{}{}
	for _,v1 := range a1{
		status := 0
		for _,v2 := range a3{
			if v1 == v2 {
				status = 1
			}
		}
		list = append(list,map[string]interface{}{"symbol":v1,"status":status})
	}

	c.JSON(200,gin.H{
		"status":1,
		"data":list,
	})
}

func DepthCommit(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	platform := c.Query("platform")
	if len(platform) == 0 {
		c.JSON(200,gin.H{"status":false})
		return
	}

	t1 := fmt.Sprintf("tmp:config:depth:%s:spot",platform)
	t2 := fmt.Sprintf("tmp:config:depth:%s:future",platform)

	if r1,_ := redis.SMembers(t1).Result(); len(r1) > 0 {
		key := fmt.Sprintf("config:depth:%s:spot",platform)
		redis.Del(key)
		redis.SAdd(key,r1)
	}

	if r2,_ := redis.SMembers(t2).Result(); len(r2) > 0 {
		key := fmt.Sprintf("config:depth:%s:future",platform)
		redis.Del(key)
		redis.SAdd(key,r2)
	}

	c.JSON(200,gin.H{"status":true})
}
