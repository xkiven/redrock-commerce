package cart

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"redrock/function/product/search"
)

type AddCartRequest struct {
	ProductId int `form:"product_id" binding:"required"`
}

func AddToCart(c *gin.Context) {
	var req AddCartRequest
	//绑定请求参数
	if err := c.ShouldBind(&req); err != nil {
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

	//添加到购物车
	_, err = db.Exec("INSERT INTO cart (user_id,product_id) VALUES (?,?)", userId, req.ProductId)
	//fmt.Println("请求的商品id:", req.ProductId)
	//fmt.Println("用户id", userId)
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
