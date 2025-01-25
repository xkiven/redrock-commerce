package user

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func Register(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//绑定请求参数
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"status": 10001,
			"info":   "参数绑定错误",
		})
		return
	}
	//连接数据库
	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "数据库连接错误",
		})
		return
	}
	defer db.Close()
	//检查用户名是否存在
	var count int

	err = db.QueryRow("select count(*) from user where username = ?", user.Username).Scan(&count)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10003,
			"info":   "查询数据库错误",
		})
		return
	}
	if count > 0 {
		c.JSON(400, gin.H{
			"status": 10004,
			"info":   "用户名已存在",
		})
		return
	}
	//插入新用户
	_, err = db.Exec("insert into user(username,password) values(?,?)", user.Username, user.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10005,
			"info":   "插入数据库失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
	})
}
