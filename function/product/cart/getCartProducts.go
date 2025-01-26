package cart

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"redrock/function/product/search"
	"strings"
)

type Product struct {
	ProductID   int    `json:"product_iD"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CommentNum  int    `json:"comment_num"`
	Price       string `json:"price"`
	IsAddedCart bool   `json:"is_addedCart"`
	Cover       string `json:"cover"`
	PublishTime string `json:"publish_time"`
	Link        string `json:"link"`
}

func GetCartProducts(c *gin.Context) {
	// 获得 token
	token := c.GetHeader("Authorization")
	if len(token) > 6 && token[:6] == "Bearer" {
		token = token[6:]
	}
	// 验证 token
	isValidToken, err := search.ValidateToken(token)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "token 验证时数据库操作出错",
		})
		return
	}
	if !isValidToken {
		c.JSON(401, gin.H{
			"status": 10003,
			"info":   "无效的 token",
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
	// 查询用户购物车
	query := "SELECT product_id FROM cart WHERE user_id=?"
	rows, err := db.Query(query, userId)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10003,
			"info":   "查询数据库错误 1",
		})
		return
	}
	defer rows.Close()

	var productIDs []interface{}
	var placeholders []string
	for rows.Next() {
		var productID int
		err := rows.Scan(&productID)
		if err != nil {
			c.JSON(500, gin.H{
				"status": 10003,
				"info":   "扫描结果错误",
			})
			return
		}
		productIDs = append(productIDs, productID)
		placeholders = append(placeholders, "?")
	}

	if len(productIDs) == 0 {
		c.JSON(200, gin.H{
			"status": 10000,
			"info":   "success",
			"data": map[string][]Product{
				"products": []Product{},
			},
		})
		return
	}

	// 第二个查询语句：根据 product_id 从 product 表中获取产品详细信息
	secondQuery := `
        SELECT 
            p.product_id, 
            p.name, 
            p.description, 
            p.type, 
            p.comment_num, 
            p.price, 
            TRUE AS is_addedCart, 
            p.cover, 
            p.publish_time, 
            p.link 
        FROM 
            product p
        JOIN 
            cart c ON p.product_id = c.product_id AND c.user_id =?
        WHERE 
            p.product_id IN (` + strings.Join(placeholders, ",") + `)`

	args := append([]interface{}{userId}, productIDs...)
	secondRows, err := db.Query(secondQuery, args...)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10003,
			"info":   "查询数据库错误 2",
		})
		return
	}
	defer secondRows.Close()

	var products []Product
	for secondRows.Next() {
		var product Product
		err := secondRows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Type,
			&product.CommentNum,
			&product.Price,
			&product.IsAddedCart,
			&product.Cover,
			&product.PublishTime,
			&product.Link)
		if err != nil {
			c.JSON(500, gin.H{
				"status": 10003,
				"info":   "扫描数据库结果错误",
			})
			return
		}
		products = append(products, product)
	}
	// 返回结果
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
		"data": map[string][]Product{
			"products": products,
		},
	})
}
