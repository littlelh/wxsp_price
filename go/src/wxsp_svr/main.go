package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "time"
    "io/ioutil"
    "encoding/json"
<<<<<<< HEAD
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    // "os"
    // "io"
    // "strings"
    // "math/rand"
=======
>>>>>>> parent of 7a2ef40... 修改http后台
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

<<<<<<< HEAD
type GoodsInfo struct {
    Id       int
    Name     string
    Describe string
    ShopType int
    Price    string
    Coupon   string
    Discount string
    ImgUrl   string
}

=======
>>>>>>> parent of 7a2ef40... 修改http后台
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
<<<<<<< HEAD
    index := c.Query("index")
    // file_name := fmt.Sprintf("%s%s%s", "/data/todd/wxsp_image/", strconv.Itoa(index), ".jpg")
    file_name := fmt.Sprintf("%s%s%s", "/data/todd/wxsp_image/", index, ".jpg")
=======
    file_name := "/data/todd/wxsp_image/2.jpg"
>>>>>>> parent of 7a2ef40... 修改http后台
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

<<<<<<< HEAD
func GetAllGoodsInfo(c *gin.Context) {
    db, err := sql.Open("mysql", "todd:temppwd@tcp(127.0.0.1:5049)/wxsp_price")
    defer db.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    db.SetMaxIdleConns(20)
    db.SetMaxOpenConns(20)

    if err := db.Ping(); err != nil{
        fmt.Println(err)
        return
    }

    rows, err := db.Query("select * from t_spider_obj;")
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return
    }

    goods_infos := make([]GoodsInfo, 0)
    for rows.Next() {
        var goods_info GoodsInfo
        var tmp_url string
        err = rows.Scan(&goods_info.Id, &goods_info.Name, &goods_info.Describe, &tmp_url, &goods_info.ShopType)
        if err != nil {
            fmt.Println(err)
            return
        }

        url := "select * from t_product_info_"
        url = fmt.Sprintf("%s%s%s", url, strconv.Itoa(goods_info.Id), " order by time_stamp DESC limit 1;")
        tmp_rows, err := db.Query(url)
        defer tmp_rows.Close()

        for tmp_rows.Next() {
            var time_str string 
            err = tmp_rows.Scan(&goods_info.Price, &goods_info.Coupon, &goods_info.Discount, &goods_info.ShopType, &time_str)
            if err != nil {
                fmt.Println(err)
            }
            
        }
        goods_infos = append(goods_infos, goods_info)
        // fmt.Println(goods_info)
    }

    c.JSON(http.StatusOK, gin.H{
        // "status":  gin.H{
        //     "goods_name":     goods_info.Name,
        // }
        "data": goods_infos,
    })
}

=======
>>>>>>> parent of 7a2ef40... 修改http后台
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
<<<<<<< HEAD
    router.GET("/allgoodsinfo", GetAllGoodsInfo)
=======
>>>>>>> parent of 7a2ef40... 修改http后台
    // router.Run()
    router.Run(":8080")
}