# **RedRock commerce**

##  **项目简介**

本项目是一个基于 Go 语言和 Gin 框架开发,借助go-cache库实现缓存机制的电商系统，提供了用户管理、商品查询、商品评论、购物车等功能。

## **技术栈**

- 编程语言：Go 语言

- Web 框架：Gin
- 数据库：MySQL
- 缓存：go - cache

## **接口说明**

### 用户相关

#### 修改用户信息：

- 接口路径： /user/update 
- 请求方法： POST 
- 请求参数：
- token ：用户身份验证令牌
- username ：新的用户名（可选）
- password ：新的密码（可选）
- 响应示例：

``` json
{
    "status": 10000,
    "info": "success"
}
```



#### 登录获取 token：

- 接口路径： /user/login 
- 请求方法： POST 
- 请求参数：
- username ：用户名
- password ：密码
- 响应示例：

```json
{
    "status": 10000,
    "info": "success",
    "token": "生成的用户 token"
}
```

#### 注册新用户：

- 接口路径： /user/register 
- 请求方法： POST 
- 请求参数：
- username ：用户名
- password ：密码
- 响应示例：

```json
{
    "status": 10000,
    "info": "success"
}
```

#### 更新token:

- 接口路径： /user/token/refresh

- 请求方法： GET

- 请求参数：

- refresh_token:更新的token

- 响应示例：

``` json
{
    "data": {
        "refresh_token": "user_1_refresh_1737824190",
        "token": "user_1_1737824190"
    },
    "info": "success",
    "status": 10000
}

```

### 商品相关

#### 搜索商品：

- 接口路径：/product/search

- 请求方法： GET 

- 请求参数：

- keyword ：商品名称关键词（可选，用于模糊查询）

响应示例：

```json
{
    "data": {
        "products": [
            {
                "product_iD": 2,
                "name": "T-shirt",
                "description": "一件短袖",
                "type": "clothes",
                "comment_num": 100,
                "price": "88.88",
                "is_addedCart": false,
                "cover": "http:/127.0.0.1/picture_url2",
                "publish_time": "1980-11-07 00:00:00",
                "link": "http://127.0.0.1/test2"
            }
        ]
    },
    "info": "success",
    "status": 10000
}
```



#### 获取商品列表：

- 接口路径： /product/list 
- 请求方法： GET 
- 请求参数：
- 响应示例：

```json
{
    "data": {
        "products": [
            {
                "product_iD": 1,
                "name": "傲慢与偏见",
                "description": "一本书",
                "type": "book",
                "comment_num": 35,
                "price": "9.80",
                "is_addedCart": false,
                "cover": "http://127.0.0.1/picture_url1",
                "publish_time": "1980-11-07 00:00:00",
                "link": "http://127.0.0.1/test1"
            },
            //...
        ]
    },
    "info": "success",
    "status": 10000
}
```



#### 查看分类下的商品：

- 接口路径： /product/:type 

- 请求方法： GET 

- 请求参数：

- type：商品分类 ID

- 响应示例：

  ```json
  {
      "status": 10000,
      "info": "success",
      "data": [
          {
              "product_iD": 1,
              "name": "傲慢与偏见",
              "description": "一本书",
              "type": "book",
              "comment_num": 35,
              "price": "9.80",
              "is_addedCart": false,
              "cover": "http://127.0.0.1/picture_url1",
              "publish_time": "1980-11-07 00:00:00",
              "link": "http://127.0.0.1/test1"
          },
          // 更多商品数据
      ]
  }
  ```



### 评论相关

#### 添加评论：

- 接口路径： /comment/:product_id
- 请求方法： POST 
- 请求参数：
- token ：用户身份验证令牌
- product_id: 商品id
- content ：评论内容
- 响应示例：

```json
{
    "data": 4,
    "info": "success",
    "status": 10000
}
```



#### 获取商品评论

- 接口路径：/comment/:product_id
- 请求方法： GET 
- 请求参数：
- product_Id ：商品 ID
- 响应示例：

```json
{
    "Status": 10000,
    "data": {
        "comments": [
            {
                "user_id": 1,
                "content": "这是一条测试评论"
            },
            {
                "user_id": 1,
                "content": "这还是一条测试评论"
            }
        ]
    },
    "info": "success"
}
```



#### 删除评论（会同时删除点赞或点踩记录）：

- 接口路径： /comment/:comment_id
- 请求方法： DELETE 
- 请求参数：
- comment_Id ：评论 ID
- 响应示例：

```json
{
    "status": 10000,
    "info": "success"
}
```

#### 点赞与点踩（对评论）：

-  接口路径：/comment/action/:comment_id

- 请求方法：POST

- 请求参数：

- comment_id: 评论 ID

- action:1为点赞，2为点踩

- 响应示例：

- ``` json
  {
      "info": "success",
      "status": 10000
  }
  ```



### 购物车相关

#### 商品加入购物车：

- 接口路径：/product/cart/addCart
- 请求方法： PUT
- 请求参数：
- token ：用户身份验证令牌
- product_Id ：商品 ID
- 响应示例：

```json
{
    "status": 10000,
    "info": "success"
}
```

#### 从购物车删除

- 接口路径：/product/cart/removeCart

- 请求方法： DELETE

- 请求参数：

- token ：用户身份验证令牌

- product_Id ：商品 ID

- 响应示例：

  ``` json
  {
      "info": "success",
      "status": 200
  }
  ```

#### 获取购物车商品列表

- 接口路径：/product/cart

- 请求方法：GET

- 请求参数：

- token ：用户身份验证令牌

- 响应示例：

  ``` json
  {
      "data": {
          "products": [
              {
                  "product_iD": 1,
                  "name": "傲慢与偏见",
                  "description": "一本书",
                  "type": "book",
                  "comment_num": 35,
                  "price": "9.80",
                  "is_addedCart": true,
                  "cover": "http://127.0.0.1/picture_url1",
                  "publish_time": "1980-11-07 00:00:00",
                  "link": "http://127.0.0.1/test1"
              }
              //...
          ]
      },
      "info": "success",
      "status": 10000
  }
  ```

### 订单相关

#### 下单

- 接口路径：/operate/order
- 请求方法：POST
- 请求参数：
- product_id：商品id
- quantity:商品数量
- address：地址
- 响应示例：

``` json
{
    "info": "success",
    "order_id": 14,
    "status": 10000
}
```

#### 删除订单

- 接口路径：/operate/order/:order_id
- 请求方法：DELETE
- 请求参数：
- token ：用户身份验证令牌
- order_id: 订单id
- 响应示例：

``` json
{
    "info": "success",
    "status": 200
}
```

#### 查看订单详情

- 接口路径：/operate/order/get/:order_id
- 请求方法：GET
- 请求参数：
- token ：用户身份验证令牌
- order_id: 订单id
- 响应示例：

``` json
{
    "data": {
        "orders": [
            {
                "user_id": 1,
                "orders": "[{\"product_id\":\"1\",\"quantity\":2},{\"product_id\":\"3\",\"quantity\":4}]",
                "address": "54Gr5pif",
                "total": 356.4
            }
        ]
    },
    "info": "success",
    "status": 10000
}
```



## 项目部署

- 确保已安装Go环境、MySQL数据库

- 通过 go get命令获取go-cache库及依赖

- 在项目根目录下执行以下命令启动服务：

```bash
go run cmd.go
```

## 特别说明

​	项目是从零开始一个人没日没夜赶在过年前完成的，缓存是做完后再加的实在是不想改了，只在搜索商品端口加了，这个时候才发现商品的is_addedCart在有无token时有变化，于是可以看到我搜索商品端口写了有token和没有token的方法，根据token得到两个结果，同理商品相关的端口都要改，但是我实在不想改了，就这样了，也还是能行


服务启动后，默认监听在  http://localhost:8080 。
