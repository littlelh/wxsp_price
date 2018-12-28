package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
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
	obj_url := "select * from t_spider_obj where pid>="
    obj_url = fmt.Sprintf("%s%s%s%s", obj_url, first_index, " and pid<=", last_index)
    rows, err := db.Query(obj_url)
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
            var time_str string 
            err = tmp_rows.Scan(&goods_info.Price, &goods_info.Coupon, &goods_info.Discount, &goods_info.ShopType, &time_str)
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