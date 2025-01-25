package search

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

// SearchWithToken 有token的搜索逻辑
func SearchWithToken(c *gin.Context, productName, token string) {
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
	err = db.QueryRow("select id from user where token = ?", token).Scan(&userId)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10004,
			"info":   "获取用户id时数据库操作出错",
		})
		return
	}
	//fmt.Println("获取到的userid:", userId)
	//fmt.Println("获取到的productName:", productName)
	query := "SELECT p.product_id, p.name, p.description, p.type, p.comment_num, p.price, c.product_id IS NOT NULL AS is_addedCart, p.cover, p.publish_time, p.link FROM product p LEFT JOIN cart c ON p.product_id = c.product_id AND c.user_id =? WHERE p.name LIKE?"
	//fullQuery := fmt.Sprintf(query, userId, "%"+productName+"%")
	//fmt.Println("执行的SQL语句:", fullQuery)
	rows, err := db.Query(query, userId, "%"+productName+"%")
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10005,
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
			&product.Link)
		if err != nil {
			c.JSON(500, gin.H{
				"status": 10003,
				"info":   "扫描数据库结果出错",
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
