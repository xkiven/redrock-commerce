package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"redrock/function/product/search"
)

// OrderItem 定义订单中的商品项结构体
type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Order 定义订单结构体
type Order struct {
	UserId  int         `json:"user_id"`
	Orders  []OrderItem `json:"orders"`
	Address interface{} `json:"address"`
	Total   float64     `json:"total"`
}

func PlaceOrder(c *gin.Context) {
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
	// 绑定参数
	var order Order
	if err := c.BindJSON(&order); err != nil {
		c.JSON(400, gin.H{
			"status": 10004,
			"info":   "请求参数绑定错误",
		})
		return
	}
	fmt.Printf("绑定后的 order: %+v\n", order)
	// 将商品项列表序列化为JSON字符串
	fmt.Printf("序列化前 order.Orders: %v\n", order.Orders)
	orderItemsJSON, err := json.Marshal(order.Orders)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10011,
			"info":   "序列化订单商品项失败",
		})
		return
	}
	fmt.Printf("序列化后 orderItemsJSON: %s\n", orderItemsJSON)

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

	// 计算总价，这里需要遍历商品项去获取每个商品价格并累加
	total := 0.0
	for _, item := range order.Orders {
		var price float64
		err = db.QueryRow("SELECT price FROM product WHERE product_id=?", item.ProductID).Scan(&price)
		if err != nil {
			c.JSON(400, gin.H{
				"status": 10003,
				"info":   "查找商品价格时数据库操作出错",
			})
			return
		}
		fmt.Printf("商品 %s 的价格: %f\n", item.ProductID, price)
		total += price * float64(item.Quantity)
	}
	fmt.Printf("计算后的总价 total: %f\n", total)
	// 开始事务
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

	// 插入订单数据到orders表，将商品项JSON字符串存储到orders字段
	var orderID int
	insertStmt := "INSERT INTO orders (user_id, orders, address, total) VALUES (?,?,?,?)"
	fmt.Printf("实际执行的SQL语句: %s, 参数: %v, %v, %v, %v\n", insertStmt, userId, string(orderItemsJSON), order.Address, total)
	result, err := tx.Exec(insertStmt, userId, string(orderItemsJSON), order.Address, total)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"status": 10007,
			"info":   "插入订单数据失败",
		})
		return
	}
	orderID64, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{
			"status": 10008,
			"info":   "获取订单id失败",
		})
		return
	}
	orderID = int(orderID64)

	// 提交事务
	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10009,
			"info":   "提交事务失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   10000,
		"info":     "success",
		"order_id": orderID,
	})
}
