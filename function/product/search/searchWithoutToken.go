package search

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
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

// SearchWithoutToken 没有token的搜索逻辑
func SearchWithoutToken(c *gin.Context, productName string) {
	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "数据库连接错误",
		})
		return
	}
	defer db.Close()

	query := fmt.Sprintf("select product_id,name,description,type,comment_num,price,is_addedCart,cover,publish_time,link from product where name like'%%%s%%'", productName)
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10003,
			"info":   "查询数据库错误",
		})
		return
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Description,
			&product.Type,
			&product.CommentNum,
			&product.Price,
			&product.IsAddedCart,
			&product.Cover,
			&product.PublishTime,
			&product.Link,
		)
		if err != nil {
			c.JSON(500, gin.H{
				"status": 10003,
				"info":   "扫描数据库结果错误",
			})
			return
		}
		products = append(products, product)
	}
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
		"data": map[string][]Product{
			"products": products,
		},
	})

}
