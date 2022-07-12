package main

import (
	_ "fmt"
	"gpn"
	"net/http"
)

func main() {
	r := gpn.New()
	//添加logger中间件
	r.Use(gpn.Logger())
	//错误恢复
	r.Use(gpn.Recovery())
	//添加组
	group := r.Group("/v1")
	r.GET("/", func(c *gpn.Context) {
		c.Data(http.StatusOK, []byte("<h1>欢迎使用gpn</h1>"))
	})
	r.GET("/lp", func(c *gpn.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	group.POST("/login", func(c *gpn.Context) {
		c.JSON(http.StatusOK, gpn.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.GET("/panic", func(c *gpn.Context) {
		names := []string{"test"}
		c.String(http.StatusOK, names[100])
	})

	group.GET("/hello/:name", func(c *gpn.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})


	r.Run(":9999")
}


