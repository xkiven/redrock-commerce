package user

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	// 从请求头中获取 Authorization 字段的值
	authHeader := c.GetHeader("Authorization")
	// 检查 Authorization 头是否存在
	if authHeader == "" {
		c.JSON(400, gin.H{
			"status": 10009,
			"info":   "Authorization头缺失",
		})
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "数据库连接问题",
		})
		return
	}
	defer db.Close()

	//查找并赋值
	var staredRefreshToken string
	var userId int
	err = db.QueryRow("select id,refresh_token from user where refresh_token=?", refreshToken).Scan(&userId, &staredRefreshToken)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 100010,
			"info":   "刷新token无效",
		})
		return
	}

	//生成新的token和refresh_token
	newToken := fmt.Sprintf("%s_%d", fmt.Sprintf("user_%d", userId), time.Now().Unix())
	newRefreshToken := fmt.Sprintf("%s_refresh_%d", fmt.Sprintf("user_%d", userId), time.Now().Unix())

	//更新数据库中的token和refresh_token
	_, err = db.Exec("update user  set token=?,refresh_token=? where id=?", newToken, newRefreshToken, userId)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10005,
			"info":   "更新数据库错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
		"data": map[string]string{
			"refresh_token": newRefreshToken,
			"token":         newToken,
		},
	})

}
