package app

import (
	"exchange_api/db"
	"exchange_api/tool"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type account struct {
	Username string	`form:"username" binding:"required"`
	Password string	`form:"password" binding:"required"`
}

func LoginIndex(c *gin.Context)  {
	c.HTML(200,"login.html",gin.H{})
}

func LoginPost(c *gin.Context)  {
	redis := db.InitRedis()
	defer redis.Close()

	var params account

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(200,gin.H{"error":err})
		return
	}

	_ = c.Request.ParseForm()
	hash,_ := redis.Get("admin:account:"+params.Username).Result()
	pass ,_ := tool.AesDecrypt(hash)
	if pass != params.Password {
		c.JSON(200,gin.H{
			"status": -1,
			"msg":"账号或密码错误",
		})
		return
	}

	generateToken(c,params)
}

// 生成令牌
func generateToken(c *gin.Context, account account ) {
	j := &tool.JWT{
		[]byte("tsewell"),
	}
	claims := tool.CustomClaims{
		account.Username,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 86400), // 过期时间 一小时
			Issuer:    "tsewell",                   //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}


	data := LoginResult{
		Username:  account.Username,
		Token: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "登录成功！",
		"data":   data,
	})
	return
}

// LoginResult 登录结果结构
type LoginResult struct {
	Token string `json:"token"`
	Username string `json:"username"`
}
