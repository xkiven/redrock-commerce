package comment

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func DeleteComment(c *gin.Context) {
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

	//开启事务(确保删除点赞记录和评论一起进行)
	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10010,
			"info":   "开启事务失败",
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//删除点赞点踩记录
	_, err = tx.Exec("delete from praise where comment_id=?", commentID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10011,
			"info":   "删除点赞或点踩记录时数据库操作出错",
		})
		return
	}

	//删除评论
	_, err = tx.Exec("delete from comment where comment_id=?", commentID)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"status": 10011,
			"info":   "删除评论时数据库操作出错",
		})
		return
	}

	//提交事务
	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10010,
			"info":   "提交事务失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
	})
}
