package main

// import (
// 	"fmt"
// )

const (
    HOUR_PRICE = "0"
	DAY_PRICE = "1"
	WEEK_PRICE = "2"
    MONTH_PRICE = "3"
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
    Price  string
    Time   string
}

func StrReverse(infos []string) []string {
	for i, j := 0, len(infos) - 1; i < j; i, j = i + 1, j - 1 {
		infos[i], infos[j] = infos[j], infos[i]
	}
	return infos
}

func FloReverse(infos []float64) []float64 {
	for i, j := 0, len(infos) - 1; i < j; i, j = i + 1, j - 1 {
		infos[i], infos[j] = infos[j], infos[i]
	}
	return infos
}

func Reverse(s []byte) []byte {
	for i, j := 0, len(s) - 1; i < j; i, j = i + 1, j - 1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}