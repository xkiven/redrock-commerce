package order

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"redrock/function/product/search"
	"strconv"
)

type orderR struct {
	UserId  int         `json:"user_id"`
	Orders  string      `json:"orders"`
	Address interface{} `json:"address"`
	Total   float64     `json:"total"`
}

func GetOrder(c *gin.Context) {
	//获取orderID
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
	fmt.Println("获取到的orderId:", orderID)
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
	//查找并处理
	query := "SELECT user_id,orders,address,total FROM orders WHERE orders_id=?"
	rows, err := db.Query(query, orderID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10003,
			"info":   "查询数据库错误",
		})
		return
	}
	defer rows.Close()
	var orders []orderR
	for rows.Next() {
		var order orderR
		err := rows.Scan(&order.UserId, &order.Orders, &order.Address, &order.Total)
		//fmt.Println("orders:", order.UserId, order.Orders, order.Address, order.Total)
		if err != nil {
			fmt.Println("Scan err:", err)
			c.JSON(500, gin.H{
				"status": 10003,
				"info":   "扫描数据库错误",
			})
			return
		}
		if userId != order.UserId {
			c.JSON(401, gin.H{
				"status": 10003,
				"info":   "用户id错误",
			})
		}
		orders = append(orders, order)
		err = rows.Err()
		if err != nil {
			c.JSON(500, gin.H{
				"status": 10005,
				"info":   "迭代查询错误",
			})
			return
		}
	}
	// 输出最终的 orders 切片长度，查看是否为空
	fmt.Printf("Orders slice length: %d\n", len(orders))
	fmt.Println("orders:", orders)
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
		"data": map[string][]orderR{
			"orders": orders,
		},
	})

}
