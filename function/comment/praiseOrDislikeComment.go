package comment

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ActionRequest struct {
	Action int `json:"action"`
}

func PraiseOrDislikeComment(c *gin.Context) {
	commentID := c.Param("comment_id")
	token := c.GetHeader("Authorization")
	if len(token) > 6 && token[:6] == "Bearer" {
		token = token[6:]
	}
	if token == "" {
		c.JSON(400, gin.H{
			"code": 10001,
			"info": "身份缺失",
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
	// 获取用户 id
	var userId int
	err = db.QueryRow("SELECT id FROM user WHERE token=?", token).Scan(&userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 10003,
			"info":   "查找用户 id 时数据库操作出错",
		})
		return
	}
	//绑定参数
	var req ActionRequest
	err = c.ShouldBindJSON(&req)
	fmt.Println("action:", req)
	if err != nil || (req.Action != 1 && req.Action != 2) {

		c.JSON(400, gin.H{
			"status": 10014,
			"info":   "绑定参数错误，1为点赞，2为点踩",
		})
		return
	}
	//检查用户是否已经对该评论点赞或点踩过了
	var count int
	err = db.QueryRow("SELECT count(*) FROM praise WHERE user_id=? AND comment_id=?", userId, commentID).Scan(&count)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10012,
			"info":   "检查点赞点踩记录时数据库操作出错",
		})
		return
	}
	if count > 0 {
		c.JSON(400, gin.H{
			"status": 10013,
			"info":   "你已对该评论进行了点赞或点踩",
		})
		return
	}
	//插入点赞点踩记录
	_, err = db.Exec("INSERT INTO praise (user_id,comment_id,model) VALUES (?,?,?)", userId, commentID, req.Action)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10011,
			"info":   "插入点赞或点踩操作时数据库操作出错",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
	})

}
