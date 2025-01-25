package user

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Login(c *gin.Context) {
	// 绑定请求参数
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		c.JSON(400, gin.H{
			"status": 10001,
			"info":   "参数不完整",
		})
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "数据库连接错误",
		})
		return
	}
	defer db.Close()
	// 检查用户名和密码
	var id int
	var storedPassword string
	err = db.QueryRow("SELECT id, password FROM user WHERE username =?", username).Scan(&id, &storedPassword)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 10006,
			"info":   "用户名不存在",
		})
		return
	}
	if storedPassword != password {
		c.JSON(400, gin.H{
			"status": 10007,
			"info":   "密码错误",
		})
		return
	}
	// 生成token（这里简单示例用用户名和时间戳拼接）
	token := fmt.Sprintf("user_%s_%d", username, time.Now().Unix())
	refreshToken := fmt.Sprintf("%s_refresh_%d", username, time.Now().Unix())
	// 更新数据库中的token和refresh_token
	_, err = db.Exec("UPDATE user SET token =?,refresh_token=? WHERE id =?", token, refreshToken, id)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10008,
			"info":   "更新token错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":        10000,
		"info":          "success",
		"token":         token,
		"refresh_token": refreshToken,
	})
}
