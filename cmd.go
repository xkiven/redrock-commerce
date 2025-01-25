package main

import (
	"github.com/gin-gonic/gin"
	"redrock/function/comment"
	"redrock/function/order"
	"redrock/function/product"
	"redrock/function/product/cart"
	"redrock/function/product/search"
	"redrock/function/user"
)

func main() {
	//http://localhost:8080
	r := gin.Default()
	//注册
	r.POST("/user/register", user.Register)
	//登录
	r.GET("/user/token", user.Login)
	//刷新token
	r.GET("/user/token/refresh", user.RefreshToken)
	//修改用户信息
	r.PUT("/user/update", user.UpdateUser)
	//获取商品列表
	r.GET("/product/list", product.GetProductList)
	//搜索商品
	r.GET("/product/search", search.SearchProduct)
	//获取相应标签的商品列表
	r.GET("/product/:type", product.GetProductByType)
	//加入购物车
	r.PUT("/product/cart/addCart", cart.AddToCart)
	//从购物车删除
	r.DELETE("/product/cart/removeCart", cart.RemoveFromCart)
	//获取购物车商品列表
	r.GET("/product/cart", cart.GetCartProducts)
	//添加评论
	r.POST("/comment/:product_id", comment.AddComment)
	//点赞和点踩评论
	r.POST("/comment/action/:comment_id", comment.PraiseOrDislikeComment)
	//删除评论
	r.DELETE("/comment/:comment_id", comment.DeleteComment)
	//获取商品评论
	r.GET("/comment/:product_id", comment.GetProductComments)
	//下单
	r.POST("/operate/order", order.PlaceOrder)
	//删除订单
	r.DELETE("/operate/order/:order_id", order.DeleteOrder)

	r.Run(":8080")
}
