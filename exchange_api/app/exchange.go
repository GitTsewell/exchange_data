package app

import (
	"exchange_api/db"
	"github.com/gin-gonic/gin"
)
var (
	okex_key = "config:exchange:okex"
	bitmex_key = "config:exchange:bitmex"
	huobi_key = "config:exchange:huobi"
	binance_key = "config:exchange:binance"
)

type exchange struct {
	Okex bool `form:"okex" json:"okex" binding:"required"`
	Bitmex bool `form:"bitmex" json:"bitmex" binding:"required"`
	Huobi bool `form:"huobi" json:"huobi" binding:"required"`
	Binance bool `form:"binance" json:"binance" binding:"required"`
}

func ExchangeEdit(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	res := exchange{
		Okex:true,
		Bitmex:true,
		Huobi:true,
		Binance:true,
	}
	if v1,_ := redis.Get(okex_key).Result();v1 == "0" {
		res.Okex = false
	}
	if v2,_ := redis.Get(bitmex_key).Result(); v2 == "0" {
		res.Bitmex = false
	}
	if v3,_ := redis.Get(huobi_key).Result(); v3 == "0" {
		res.Huobi = false
	}
	if v4,_ := redis.Get(binance_key).Result(); v4 == "0" {
		res.Binance = false
	}

	data := map[string]interface{}{
		"okex":res.Okex,
		"bitmex":res.Bitmex,
		"huobi":res.Huobi,
		"binance":res.Binance,
	}

	c.JSON(200,gin.H{"data":data})
}

func ExchangeUpdate(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	var exchange exchange
	_ = c.ShouldBindJSON(&exchange)

	redis.Set(okex_key,exchange.Okex,0)
	redis.Set(bitmex_key,exchange.Bitmex,0)
	redis.Set(huobi_key,exchange.Huobi,0)
	redis.Set(binance_key,exchange.Binance,0)

	c.JSON(200,gin.H{"status":1})
}
