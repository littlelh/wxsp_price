package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const (
    HOUR_PRICE = "0"
    DAY_PRICE = "1"
    MONTH_PRICE = "2"
)

var (
	db    *sql.DB
)

func init_mysql() bool {
	var err error
	db, err = sql.Open("mysql", "todd:temppwd@tcp(127.0.0.1:5049)/wxsp_price")
    if err != nil {
        fmt.Println("db open failed,  err:", err)
        return false
    }

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		fmt.Println("db ping failed, err:", err)
		return false
	}

	fmt.Println("mysql module init success.")
	return true
}

func QueryGoodsInfo(goods_infos *[]GoodsInfo, first_index string, last_index string) bool {
	obj_sql := "select * from t_spider_obj where pid>="
    obj_sql = fmt.Sprintf("%s%s%s%s", obj_sql, first_index, " and pid<=", last_index)
    rows, err := db.Query(obj_sql)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return false
	}

	for rows.Next() {
        var goods_info GoodsInfo
        var tmp_url string
        err = rows.Scan(&goods_info.Id, &goods_info.Name, &goods_info.Describe, &tmp_url, &goods_info.ShopType)
        if err != nil {
            fmt.Println(err)
            return false
        }

        url := "select * from t_product_info_"
        url = fmt.Sprintf("%s%s%s", url, strconv.Itoa(goods_info.Id), " order by time_stamp DESC limit 1")
        tmp_rows, err := db.Query(url)
        defer tmp_rows.Close()

        for tmp_rows.Next() {
            err = tmp_rows.Scan(&goods_info.Price, &goods_info.Coupon, &goods_info.Discount, &goods_info.ShopType, &goods_info.UpdateTime)
            if err != nil {
				fmt.Println(err)
				return false
            }
		}
		// fmt.Println(goods_info)
        *goods_infos = append(*goods_infos, goods_info)
	}
	return true
}

func QueryAllGoodsInfo(goods_info []GoodsInfo) {
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
}

func QueryHourGoodsPrice(price_infos *[]PriceInfo, goods_id string) bool {
    sql := "select price,coupon,time_stamp from t_product_info_" + goods_id + " order by time_stamp DESC limit 12"
    rows, err := db.Query(sql)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return false
	}

    for rows.Next() {
        var price_info PriceInfo
        var str_tmp_price, str_tmp_coupon string
        var tmp_price, tmp_coupon float64
        err = rows.Scan(&str_tmp_price, &str_tmp_coupon, &price_info.Time)
        if err != nil {
            fmt.Println(err)
            return false
        }

        price_info.Time = string([]rune(price_info.Time)[11:])

        tmp_price, err = strconv.ParseFloat(str_tmp_price, 32)
        if err != nil {
            fmt.Println(err)
            return false
        }

        tmp_coupon, err = strconv.ParseFloat(str_tmp_coupon, 32)
        if err != nil {
            fmt.Println(err)
            return false
        }

        price_info.Price = strconv.FormatFloat(tmp_price, 'f', 2, 64)
        price_info.Coupon = strconv.FormatFloat(tmp_coupon, 'f', 2, 64)
        
        *price_infos = append(*price_infos, price_info)
    }
    reverse(*price_infos)
    return true
}

func reverse(s []PriceInfo) []PriceInfo {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}