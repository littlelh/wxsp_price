package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "time"
    "io/ioutil"
    "encoding/json"
    "os"
    // "io"
    // "strings"
    // "math/rand"
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

type GoodsInfo struct {
    Id       int
    Name     string
    Describe string
    ShopType int
    Price    string
    Coupon   string
    Discount string
    ImgUrl   string
    UpdateTime string
}

type PriceInfo struct {
    Price string
    Time  string
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
    })
}

func LoadImage(c *gin.Context) {
    index := c.Query("index")
    // file_name := fmt.Sprintf("%s%s%s", "/data/todd/wxsp_image/", strconv.Itoa(index), ".jpg")
    file_name := fmt.Sprintf("%s%s%s", "/data/todd/wxsp_image/", index, ".jpg")
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

func GetGoodsInfo(c *gin.Context) {
    first_index := c.DefaultQuery("first", "0")
    last_index := c.DefaultQuery("last", "0")

    if first_index == "0" || last_index == "0" {
        fmt.Println("param failed")
        return
    }

    goods_infos := make([]GoodsInfo, 0)
    if !QueryGoodsInfo(&goods_infos, first_index, last_index) {
        fmt.Println("query goods info from mysql failed")
        return
    }

    c.JSON(http.StatusOK, gin.H{
        // "status":  gin.H{
        //     "goods_name":     goods_info.Name,
        // }
        "data": goods_infos,
    })
}

func GetGoodsPrice(c *gin.Context) {
    goods_id := c.DefaultQuery("goods_id", "0")
    query_type := c.DefaultQuery("query_type", "0")

    price_infos := make([]PriceInfo, 0)

    switch {
        case query_type == HOUR_PRICE:
            if !QueryHourGoodsPrice(&price_infos, goods_id) {
                fmt.Println("query price info from mysql failed")
            }
        case query_type == DAY_PRICE:
        case query_type == MONTH_PRICE:
        default:
            fmt.Println("not support this query type")
    }

    c.JSON(http.StatusOK, gin.H{
        "data": price_infos,
    })
}

func main() {
    // gin.DisableConsoleColor()
    // log_file, _ := os.Create("gin.log")
    // gin.DefaultWriter = io.MultiWriter(log_file)

    ok := init_mysql()
	if !ok {
		fmt.Println("init mysql failed, exit")
		os.Exit(1)
	}

    router := gin.Default()
    // router_v1 := router.Group("/v1")
    // 用户登录
    router.GET("/login", UserLogin)

    // 下载图片
    router.GET("/loadimage", LoadImage)

    // 获取所有商品信息
    router.GET("/goodsinfo", GetGoodsInfo)

    // 获取价格波动信息
    router.GET("/goodsprice", GetGoodsPrice)

    // router.Run()
    router.Run(":8080")

    db.Close()
}