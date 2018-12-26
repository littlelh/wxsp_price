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
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "math/rand"
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
    num := rand.Intn(8) + 1
    fmt.Println(num)
    file_name := fmt.Sprintf("%s%s%s", "/data/todd/wxsp_image/", strconv.Itoa(num), ".jpg")
    // file_name := "/data/todd/wxsp_image/2.jpg"
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

func GetAllGoodsInfo(c *gin.Context) {
    db, err := sql.Open("mysql", "todd:temppwd@tcp(127.0.0.1:5049)/wxsp_price")
    defer db.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    rows, err := db.Query("select * from t_spider_obj;")
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    for rows.Next() {
        var good_info GoodsInfo
        var tmp_url string
        err = rows.Scan(&good_info.Id, &good_info.Name, &good_info.Describe, &tmp_url, &good_info.ShopType)
        if err != nil {
            fmt.Println(err)
            return
        }

        url := "select * from t_product_info_"
        url = fmt.Sprintf("%s%s%s", url, strconv.Itoa(good_info.Id), " order by time_stamp DESC limit 1;")
        tmp_rows, err := db.Query(url)
        defer tmp_rows.Close()

        for tmp_rows.Next() {
            var time_str string 
            err = tmp_rows.Scan(&good_info.Price, &good_info.Coupon, &good_info.Discount, &good_info.ShopType, &time_str)
            if err != nil {
                fmt.Println(err)
            }
            
        }
        // fmt.Println(good_info)
    }

    // c.JSON(http.StatusOK, gin.H{
    //     "status":  gin.H{
    //         "status_code": http.StatusOK,
    //         "status":      "ok",
    //     },
    //     "message": message,
    //     "nick":    nick,
    // })

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
    router.GET("/goodsinfo", GetAllGoodsInfo)
    // router.Run()
    router.Run(":8080")
}