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

func QueryHourGoodsPrice(price_infos *[]float64, time_infos *[]string, goods_id string) bool {
    sql := "select CAST(price as DECIMAL) price, CAST(coupon as DECIMAL) coupon, time_stamp from t_product_info_" +
        goods_id + " order by time_stamp DESC limit 12"
    rows, err := db.Query(sql)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return false
	}

    for rows.Next() {
        var tmp_time string
        var tmp_price, tmp_coupon float64
        err = rows.Scan(&tmp_price, &tmp_coupon, &tmp_time)
        if err != nil {
            fmt.Println(err)
            return false
        }
        // tmp_price, err = strconv.ParseFloat(str_tmp_price, 32)
        // if err != nil {
        //     fmt.Println(err)
        //     return false
        // }
        // price_info.Time = string([]rune(price_info.Time)[11:])
        // price_info.Price = strconv.FormatFloat(tmp_price - tmp_coupon, 'f', 2, 64)
        tmp_time = string([]rune(tmp_time)[11:])

        *price_infos = append(*price_infos, tmp_price - tmp_coupon)
        *time_infos = append(*time_infos, tmp_time)
    }
    FloReverse(*price_infos)
    StrReverse(*time_infos)
    return true
}

func QueryDayGoodsPrice(price_infos *[]float64, time_infos *[]string, goods_id string) bool {
    sql := "select min(CAST(price as DECIMAL)-CAST(coupon as DECIMAL)) min_price, DATE_FORMAT(time_stamp, '%Y%m%d') day from t_product_info_" +
        goods_id + " group by day order by day desc limit 12"
    rows, err := db.Query(sql)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return false
    }

    for rows.Next() {
        var tmp_time string
        var tmp_price float64
        err = rows.Scan(&tmp_price, &tmp_time)
        if err != nil {
            fmt.Println(err)
            return false
        }

        *price_infos = append(*price_infos, tmp_price)
        *time_infos = append(*time_infos, tmp_time)
    }
    FloReverse(*price_infos)
    StrReverse(*time_infos)
    return true
}

func QueryWeekGoodsPrice(price_infos *[]float64, time_infos *[]string, goods_id string) bool {
    sql := "select min(CAST(price as DECIMAL)-CAST(coupon as DECIMAL)) min_price, DATE_FORMAT(time_stamp, '%Y%u') week from t_product_info_" +
        goods_id + " group by week order by week desc limit 12"
    rows, err := db.Query(sql)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return false
    }

    for rows.Next() {
        var tmp_time string
        var tmp_price float64
        err = rows.Scan(&tmp_price, &tmp_time)
        if err != nil {
            fmt.Println(err)
            return false
        }

        *price_infos = append(*price_infos, tmp_price)
        *time_infos = append(*time_infos, tmp_time)
    }
    FloReverse(*price_infos)
    StrReverse(*time_infos)
    return true
}

func QueryMonthGoodsPrice(price_infos *[]float64, time_infos *[]string, goods_id string) bool {
    sql := "select min(CAST(price as DECIMAL)-CAST(coupon as DECIMAL)) min_price, DATE_FORMAT(time_stamp, '%Y%m') month from t_product_info_" +
        goods_id + " group by month order by month desc limit 12"
    rows, err := db.Query(sql)
    defer rows.Close()
    if err != nil {
        fmt.Println(err)
        return false
    }

    for rows.Next() {
        var tmp_time string
        var tmp_price float64
        err = rows.Scan(&tmp_price, &tmp_time)
        if err != nil {
            fmt.Println(err)
            return false
        }
        
        *price_infos = append(*price_infos, tmp_price)
        *time_infos = append(*time_infos, tmp_time)
    }
    FloReverse(*price_infos)
    StrReverse(*time_infos)
    return true
}
