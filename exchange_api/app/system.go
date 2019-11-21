package app

import (
	"exchange_api/db"
	"exchange_api/tool"
	"github.com/gin-gonic/gin"
)

func SystemIndex(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	k1 := "config:system:restart_local_ws"
	k2 := "config:system:restart_other_ws"

	v1,_ := redis.Get(k1).Result()
	v2,_ := redis.Get(k2).Result()

	data := []map[string]interface{}{
		{"name":"重启本机ws客户端","key":k1,"action":v1,"show":false},{"name":"重启其他服务器WS客户端","key":k2,"action":v2,"show":false},
	}

	c.JSON(200,gin.H{
		"data":data,
	})
}

func SystemExec(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	key := c.Param("key")

	action , _ := redis.Get(key).Result()

	tool.SystemExec(action)

	c.JSON(200,gin.H{"status":1})
}

func SystemUpdate(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	var data struct{
		Key string `json:"key" binding:"required"`
		Action string `json:"action" binding:"required"`
	}

	_ = c.BindJSON(&data)

	if a,_ :=redis.Set(data.Key,data.Action,0).Result();a == "OK" {
		c.JSON(200,gin.H{"status":1})
	}else {
		c.JSON(200,gin.H{"status":-1})
	}
}

func SystemProcess(c *gin.Context)  {

}
