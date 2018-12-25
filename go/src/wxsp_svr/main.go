package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    // "os"
    // "io"
    // "strings"
    "time"
    "io/ioutil"
    "encoding/json"
)

type AppInfo struct {
    Appid      string `form:"appid" json:"appid" binding:"required"`
    Secret     string `form:"secret" json:"secret" binding:"required"`
    Js_code    string `form:"js_code" json:"js_code binding:"required"`
    Grant_type string `form:"grant_type" json:"grant_type binding:"required"`
}

type UserInfo struct {
    SessionKey string `json:"session_key"`
    OpenId     string `json:"openid"`
}

func GetOpenIdAndSessionKey(app_info AppInfo) (user_info UserInfo) {
    url := "https://api.weixin.qq.com/sns/jscode2session?" + "appid=" + app_info.Appid + "&secret=" +
        app_info.Secret + "&js_code=" + app_info.Js_code + "&grant_type=" + app_info.Grant_type
    req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("http NewRequest failed, err: ", err)
		return
    }

    client := &http.Client {
		Timeout: time.Duration(5) * time.Second,
    }

    resp, err := client.Do(req)
    defer resp.Body.Close()
	if err != nil {
		fmt.Println("http failed, err: ", err)
		return
    }

    body, _ := ioutil.ReadAll(resp.Body)

    var info UserInfo
    if err := json.Unmarshal(body, &info); err != nil {
        fmt.Println("Unmarshal json error:", err)
    }

    return info
}

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

func UserLogin(c *gin.Context) {
    var app_info AppInfo

    err := c.Bind(&app_info)
    if err != nil {
        fmt.Println(err)
        return
    }

    user_info := GetOpenIdAndSessionKey(app_info)
    fmt.Println(user_info)

    c.JSON(http.StatusOK, gin.H {
        "openid": user_info.OpenId,
        // "session_key": "222",
    })
}

func LoadImage(c *gin.Context) {
    file_name := "/data/todd/wxsp_image/2.jpg"
    file, err := ioutil.ReadFile(file_name)
    if err != nil {
        fmt.Println("no such picture:", file_name)
        return
    }

    // c.Header("Content-Type", "image/png")
    // c.Header("Content-Disposition", "inline")
    c.Header("Content-Disposition", `attachment; filename=` + file_name)
    c.Data(http.StatusOK, "multipart/form-data", file)
}

func main() {
    // gin.DisableConsoleColor()
    // log_file, _ := os.Create("gin.log")
    // gin.DefaultWriter = io.MultiWriter(log_file)

    router := gin.Default()
    // router_v1 := router.Group("/v1")
    router.GET("/someGet", getting)
    router.GET("/func1", func1)
    router.GET("/login", UserLogin)
    router.GET("/loadimage", LoadImage)
    // router.Run()
    router.Run(":8080")
}