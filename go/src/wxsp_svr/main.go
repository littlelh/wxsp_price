package main

import (
    // "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    // "os"
    // "io"
)

func getting(c *gin.Context) {
    // fmt.Fprintln(gin.DefaultWriter, "foo bar")
    c.String(http.StatusOK, "Hello world")
}

// 获取请求参数的例子
func func1(c *gin.Context) {
    // DefaultQuery只作用于key不存在的时候，提供默认值
    firstname := c.DefaultQuery("firstname", "Guest")
    lastname := c.Query("lastname")
 
    c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
}

func main() {
    // gin.DisableConsoleColor()
    // log_file, _ := os.Create("gin.log")
    // gin.DefaultWriter = io.MultiWriter(log_file)

    router := gin.Default()
    router_v1 := router.Group("/v1")
    router_v1.GET("/someGet", getting)
    router_v1.GET("/func1", func1)
    // router.Run()
    router.Run(":8080")
}