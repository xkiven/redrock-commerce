package user

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"strings"
)

func UpdateUser(c *gin.Context) {
	var user struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//绑定请求参数
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"status": 10001,
			"info":   "参数绑定错误",
		})
		return
	}
	//检查token是否存在
	if user.Token == "" {
		c.JSON(400, gin.H{
			"status": 10001,
			"info":   "token参数缺失",
		})
	}
	//验证token有效性，从数据库中查询用户信息
	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "数据库连接错误",
		})
		return
	}
	defer db.Close()
	var userId int
	err = db.QueryRow("select id from user where token=?", user.Token).Scan(&userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 10011,
			"info":   "token无效或用户未登录",
		})
		return
	}
	//更新用户信息
	var updateStmt strings.Builder
	updateStmt.WriteString("UPDATE user SET ")
	updateParams := []interface{}{}
	hasField := false
	if user.Username != "" {
		if hasField {
			updateStmt.WriteString(", ")
		}
		updateStmt.WriteString("username =?")
		updateParams = append(updateParams, user.Username)
		hasField = true
	}
	if user.Password != "" {
		if hasField {
			updateStmt.WriteString(", ")
		}
		updateStmt.WriteString("password =?")
		updateParams = append(updateParams, user.Password)
		hasField = true
	}
	if !hasField {
		c.JSON(400, gin.H{
			"status": 10004,
			"info":   "没有需要更新的用户信息字段",
		})
		return
	}
	updateStmt.WriteString(" WHERE id =?")
	updateParams = append(updateParams, userId)

	_, err = db.Exec(updateStmt.String(), updateParams...)
	if err != nil {
		// 打印详细错误信息到控制台，方便排查问题
		println("数据库执行错误:", err.Error())
		c.JSON(500, gin.H{
			"status": 10005,
			"info":   "更新数据库错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
	})

}
