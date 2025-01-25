package order

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"redrock/function/product/search"
	"strconv"
)

func DeleteOrder(c *gin.Context) {
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
	// 验证token
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
	//获取想要删除的订单
	orderIDStr := c.Param("order_id")
	if len(orderIDStr) == 0 {
		c.JSON(400, gin.H{
			"status": 10013,
			"info":   "缺失订单id",
		})
		return
	}
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 10014,
			"info":   "订单id参数格式问题",
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

	//删除订单数据
	deleteStmt := "delete from orders where orders_id=?"

	_, err = db.Exec(deleteStmt, orderID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10015,
			"info":   "删除订单数据失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
	})
}
