package comment

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"redrock/function/product/search"
	"strconv"
)

type Comments struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Content   string `json:"content"`
}

func AddComment(c *gin.Context) {
	productIDStr := c.Param("product_id")

	// 将字符串类型的 product_id 转换为 int 类型
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		// 如果转换失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的产品 ID，必须是整数",
		})
		return
	}
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
	//验证token
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
	//绑定comment参数
	var comment Comments
	err = c.BindJSON(&comment)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 10005,
			"info":   "参数绑定错误",
		})
		return
	}
	//fmt.Println("comment.projectId:", comment.ProductID)
	//fmt.Println("productID:", productID)
	if comment.ProductID != productID {
		c.JSON(400, gin.H{
			"status": 10006,
			"info":   "商品id不一致",
		})
		return
	}
	if comment.Content == "" {
		c.JSON(400, gin.H{
			"status": 10007,
			"info":   "评论内容不能为空",
		})
		return
	}
	//插入评论
	result, err := db.Exec("INSERT INTO comment(user_id,product_id,content) VALUES(?,?,?)", userId, productID, comment.Content)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10009,
			"info":   "插入评论时数据库操作出错",
		})
		return
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{
			"status": 10009,
			"info":   "获取评论id时出现问题",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 10000,
		"info":   "success",
		"data":   commentID,
	})

}
