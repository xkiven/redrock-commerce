package search

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
)

var c *cache.Cache

func init() {
	//创建缓存实例，设置默认过期时间和清理时间
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func SearchProduct(c *gin.Context) {
	productName := c.Query("product_name")
	if productName == "" {
		c.JSON(400, gin.H{
			"status": 10001,
			"info":   "商品名称参数product_name是必填项",
		})
		return
	}
	token := c.GetHeader("Authorization")
	if len(token) > 6 && token[:6] == "Bearer" {
		token = token[6:]
	}
	if token == "" {
		//没有token,is_addedCart设为false
		SearchWithoutToken(c, productName)
		return
	}
	//fmt.Println("获取到的token:", token)
	//验证token
	isValidToken, err := ValidateToken(token)
	//fmt.Println("获取到的isValidToken：", isValidToken)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "token验证时数据库操作出错",
		})
		return
	}
	if !isValidToken {
		c.JSON(401, gin.H{
			"status": 10003,
			"info":   "无效的token",
		})
		return
	}
	//token验证通过，查询数据库判断商品是否在购物车
	SearchWithToken(c, productName, token)

}
