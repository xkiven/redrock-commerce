package cart

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"redrock/function/product/search"
)

type RemoveRequest struct {
	ProductId int `form:"product_id" binding:"required"`
}

func RemoveFromCart(c *gin.Context) {
	var req RemoveRequest
	//绑定请求参数
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{
			"status": 10001,
			"info":   "参数验证失败",
		})
		return
	}
	//验证token
	token := c.GetHeader("Authorization")
	if len(token) > 6 && token[:6] == "Bearer" {
		token = token[6:]
	}
	isValidToken, err := search.ValidateToken(token)
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

	//获取用户id
	var userId int
	err = db.QueryRow("SELECT id FROM user WHERE token=?", token).Scan(&userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 10003,
			"info":   "查找用户id时数据库操作出错",
		})
		return
	}
	// 查询购物车中是否存在对应记录
	var count int
	query := "SELECT COUNT(*) FROM cart WHERE user_id=? AND product_id=?"
	err = db.QueryRow(query, userId, req.ProductId).Scan(&count)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10004,
			"info":   "查询购物车记录数量时数据库操作出错",
		})
		return
	}
	if count == 0 {
		c.JSON(200, gin.H{
			"status": 200,
			"info":   "购物车中已删除",
		})
		return
	}
	if count > 0 {
		//从购物车中删除
		_, err = db.Exec("DELETE FROM cart WHERE user_id= ? AND product_id= ? ", userId, req.ProductId)
		if err != nil {
			c.JSON(500, gin.H{
				"status": 10005,
				"info":   "删除数据库失败",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
	})
}
