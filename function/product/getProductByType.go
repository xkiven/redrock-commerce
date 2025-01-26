package product

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func GetProductByType(c *gin.Context) {
	//从Path获取参数type
	productType := c.Param("type")
	//productType := c.Query("type")
	//fmt.Println("获取到的productType:", productType)
	db, err := sql.Open("mysql", "root:xkw510724@tcp(127.0.0.1:3306)/redrock_ecommerce?charset=utf8")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10002,
			"info":   "数据库连接错误",
		})
		return
	}
	defer db.Close()

	var products []Product
	query := "SELECT product_id,name,description,type,comment_num,price,is_addedCart,cover,publish_time,link FROM product WHERE type =?"
	rows, err := db.Query(query, productType)
	if err != nil {
		c.JSON(500, gin.H{
			"Status": 10003,
			"Info":   "查询数据库错误",
		})
		return
	}
	defer rows.Close()
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
			&product.Link)
		if err != nil {
			c.JSON(500, gin.H{
				"Status": 10003,
				"Info":   "扫描数据库错误",
			})
			return
		}
		products = append(products, product)
		err = rows.Err()
		if err != nil {
			c.JSON(500, gin.H{
				"Status": 10005,
				"info":   "迭代查询结果错误",
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
		"data":   products,
	})
}
