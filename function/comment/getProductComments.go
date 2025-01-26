package comment

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

type commentR struct {
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

func GetProductComments(c *gin.Context) {
	productID := c.Param("product_id")
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

	var comments []commentR
	//查找并处理
	query := "SELECT user_id,content FROM comment WHERE product_id = ?"
	rows, err := db.Query(query, productID)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": 10003,
			"Info":   "查询数据库错误",
		})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var comment commentR
		err = rows.Scan(&comment.UserID, &comment.Content)

		if err != nil {
			c.JSON(500, gin.H{
				"Status": 10003,
				"Info":   "扫描数据库错误",
			})
			return
		}
		comments = append(comments, comment)
		err = rows.Err()
		if err != nil {
			c.JSON(500, gin.H{
				"Status": 10005,
				"info":   "迭代查询错误",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"Status": 10000,
		"info":   "success",
		"data": map[string][]commentR{
			"comments": comments,
		},
	})

}
